package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sayskar500/lumina-plane/internal/models"
)

type Provider interface {
	Generate(ctx context.Context, prompt string) (*models.AIResponse, error)
	Name() string
}

type GroqProvider struct {
	APIKey string
	Model  string
}

func NewGroqProvider(apiKey, model string) *GroqProvider {
	return &GroqProvider{APIKey: apiKey, Model: model}
}

func (g *GroqProvider) Name() string { return "Groq" }

func (g *GroqProvider) Generate(ctx context.Context, prompt string) (*models.AIResponse, error) {
	if g.APIKey == "" {
		return nil, fmt.Errorf("Groq API Key is missing. Please set GROQ_API_KEY environment variable")
	}

	start := time.Now()
	url := "https://api.groq.com/openai/v1/chat/completions"

	requestBody := map[string]interface{}{
		"model": g.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+g.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error connecting to Groq: %v", err)
	}
	defer resp.Body.Close()

	// Handle specific HTTP errors for better debugging
	switch resp.StatusCode {
	case http.StatusOK:
		// Continue to decode
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("invalid Groq API Key (401 Unauthorized)")
	case http.StatusForbidden:
		return nil, fmt.Errorf("Groq API access forbidden (403 Forbidden)")
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("Groq API rate limit exceeded (429 Too Many Requests)")
	case http.StatusNotFound:
		return nil, fmt.Errorf("Groq model %s not found (404 Not Found)", g.Model)
	default:
		return nil, fmt.Errorf("unexpected Groq API response: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Groq response: %v", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("Groq returned an empty response")
	}

	return &models.AIResponse{
		Text:       result.Choices[0].Message.Content,
		TokensUsed: result.Usage.TotalTokens,
		ModelUsed:  g.Model,
		LatencyMs:  time.Since(start).Milliseconds(),
	}, nil
}
