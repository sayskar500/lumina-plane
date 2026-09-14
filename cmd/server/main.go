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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
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
}

// initOTel wires the OpenTelemetry SDK to an OTLP/gRPC collector when enabled
// in config. When disabled, it returns a no-op shutdown and the global tracer
// provider stays a no-op, so tracing calls remain essentially free.
func initOTel(ctx context.Context, cfg *config.Config) (func(context.Context), error) {
	if !cfg.Observability.OTelEnabled || cfg.Observability.OTelEndpoint == "" {
		return func(context.Context) {}, nil
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Observability.OTelEndpoint),
		otlptracegrpc.WithInsecure(), // internal collector traffic; add TLS when the collector requires it
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			"",
			attribute.String("service.name", cfg.Observability.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("⚠️ OpenTelemetry shutdown: %v", err)
		}
	}, nil
}

func main() {
	ctx := context.Background()

	var err error
	appCfg, err = config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	shutdown, err := initOTel(ctx, appCfg)
	if err != nil {
		log.Printf("⚠️ OpenTelemetry disabled (%v)", err)
	}
	defer shutdown(context.Background())
	tracer = otel.Tracer(appCfg.Observability.ServiceName)

	appSecrets, err = config.LoadSecrets("config/secrets.yaml")
	if err != nil {
		log.Printf("⚠️ Warning: Could not load secrets.yaml: %v", err)
	}
	if appSecrets.GroqAPIKey == "" {
		log.Printf("⚠️ Warning: No Groq API key found (config/secrets.yaml or GROQ_API_KEY). AI features will be disabled.")
	} else {
		aiProvider = gateway.NewGroqProvider(appSecrets.GroqAPIKey, appCfg.AI.DefaultModel)
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

	store = db.NewStore(mongoClient, appCfg.MongoDB.Database)

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
		http.Error(w, "AI Provider not configured. Check config/secrets.yaml or GROQ_API_KEY", http.StatusInternalServerError)
		return
	}

	var req models.AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Honor the caller's requested model; fall back to the configured default.
	model := req.Model
	if model == "" {
		model = appCfg.AI.DefaultModel
	}

	finalPrompt := req.Prompt
	if prompt, err := store.GetActivePrompt(ctx, req.ProjectID); err == nil {
		finalPrompt = fmt.Sprintf("%s\n\nUser Input: %s", prompt.Template, req.Prompt)
	}

	genCtx, genSpan := tracer.Start(ctx, "ai_generation")
	genStart := time.Now()
	resp, err := aiProvider.Generate(genCtx, finalPrompt, model)
	genSpan.End()

	if err != nil {
		aiRequestsTotal.WithLabelValues(model, "error").Inc()
		http.Error(w, fmt.Sprintf("AI Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	aiRequestsTotal.WithLabelValues(resp.ModelUsed, "success").Inc()
	tokenUsage.WithLabelValues(resp.ModelUsed).Add(float64(resp.TokensUsed))
	requestLatency.WithLabelValues(resp.ModelUsed).Observe(time.Since(genStart).Seconds())

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

	if p.ProjectID == "" || p.Template == "" {
		http.Error(w, "project_id and template are required", http.StatusBadRequest)
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
