package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sayskar500/lumina-plane/internal/config"
	"github.com/sayskar500/lumina-plane/internal/db"
	"github.com/sayskar500/lumina-plane/internal/gateway"
	"github.com/sayskar500/lumina-plane/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.opentelemetry.io/otel/trace"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
}

var (
	mongoClient *mongo.Client
	store       *db.Store
	aiProvider  gateway.Provider
	tracer      trace.Tracer
	appCfg      *config.Config
	appSecrets  *config.Secrets

	aiRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lumina_ai_requests_total",
			Help: "Total number of AI requests processed",
		},
		[]string{"model", "status"},
	)
	tokenUsage = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lumina_tokens_used_total",
			Help: "Total number of tokens consumed",
		},
		[]string{"model"},
	)
	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "lumina_ai_request_duration_seconds",
			Help:    "Latency of AI requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"model"},
	)
)

func init() {
	prometheus.MustRegister(aiRequestsTotal)
	prometheus.MustRegister(tokenUsage)
	prometheus.MustRegister(requestLatency)
	tracer = trace.NewNoopTracerProvider().Tracer("lumina-plane-server")
}

func initOTel(ctx context.Context) func(context.Context) {
	return func(ctx context.Context) {}
}

func main() {
	ctx := context.Background()
	shutdown := initOTel(ctx)
	defer shutdown(ctx)

	// Load Configurations from YAML
	var err error
	appCfg, err = config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	appSecrets, err = config.LoadSecrets("config/secrets.yaml")
	if err != nil {
		log.Printf("⚠️ Warning: Could not load secrets.yaml: %v. AI features will be disabled.", err)
	}

	// MongoDB Connection with Retry Logic
	for i := 1; i <= 5; i++ {
		fmt.Printf("Connecting to MongoDB (Attempt %d/5)... \n", i)
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		mongoClient, err = mongo.Connect(ctxTimeout, options.Client().ApplyURI(appCfg.MongoDB.URI))
		if err == nil {
			if err = mongoClient.Ping(ctxTimeout, readpref.Primary()); err == nil {
				fmt.Println("✅ Connected successfully to MongoDB")
				cancel()
				break
			}
		}
		fmt.Printf("⚠️ MongoDB not ready yet: %v. Retrying in 2s...\n", err)
		cancel()
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ Could not connect to MongoDB after 5 attempts:", err)
	}

	store = db.NewStore(mongoClient)
	if appSecrets != nil && appSecrets.GroqAPIKey != "" {
		aiProvider = gateway.NewGroqProvider(appSecrets.GroqAPIKey, appCfg.AI.DefaultModel)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ask", askHandler)
	http.HandleFunc("/prompt", promptHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Lumina-Plane Server starting on :%d...\n", appCfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", appCfg.Server.Port), nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := mongoClient.Ping(ctx, readpref.Primary())
	status, dbStatus := "healthy", "connected"
	if err != nil {
		status, dbStatus = "unhealthy", fmt.Sprintf("disconnected: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(HealthResponse{Status: status, Database: dbStatus})
}

func askHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "askHandler")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if aiProvider == nil {
		http.Error(w, "AI Provider not configured. Check config/secrets.yaml", http.StatusInternalServerError)
		return
	}

	var req models.AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	finalPrompt := req.Prompt
	if prompt, err := store.GetActivePrompt(ctx, req.ProjectID); err == nil {
		finalPrompt = fmt.Sprintf("%s\n\nUser Input: %s", prompt.Template, req.Prompt)
	}

	genCtx, genSpan := tracer.Start(ctx, "ai_generation")
	resp, err := aiProvider.Generate(genCtx, finalPrompt)
	genSpan.End()

	if err != nil {
		aiRequestsTotal.WithLabelValues(req.Model, "error").Inc()
		http.Error(w, fmt.Sprintf("AI Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	aiRequestsTotal.WithLabelValues(resp.ModelUsed, "success").Inc()
	tokenUsage.WithLabelValues(resp.ModelUsed).Add(float64(resp.TokensUsed))
	requestLatency.WithLabelValues(resp.ModelUsed).Observe(time.Since(time.Now()).Seconds())

	go func() {
		logCtx, logCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer logCancel()
		store.LogTokens(logCtx, models.TokenLog{
			ProjectID: req.ProjectID,
			Tokens:    resp.TokensUsed,
			Timestamp: time.Now(),
			Model:     resp.ModelUsed,
		})
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func promptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p models.Prompt
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := store.SavePrompt(ctx, p); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save prompt: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Prompt saved successfully for project %s", p.ProjectID)
}
