package models

import "time"

// AIRequest represents an incoming request to the AI Gateway
type AIRequest struct {
	ProjectID string `json:"project_id" bson:"project_id"`
	Prompt    string `json:"prompt" bson:"prompt"`
	Model     string `json:"model" bson:"model"`
}

// AIResponse represents the response from the AI Gateway
type AIResponse struct {
	Text       string `json:"text"`
	TokensUsed int    `json:"tokens_used"`
	ModelUsed  string `json:"model_used"`
	LatencyMs  int64  `json:"latency_ms"`
}

// TokenLog stores the usage data for billing and optimization
type TokenLog struct {
	ID        string    `bson:"_id,omitempty"`
	ProjectID string    `bson:"project_id"`
	Tokens    int       `bson:"tokens"`
	Timestamp time.Time `bson:"timestamp"`
	Model     string    `bson:"model"`
}

// Prompt represents a versioned prompt template
type Prompt struct {
	ID        string    `bson:"_id,omitempty"`
	ProjectID string    `json:"project_id" bson:"project_id"`
	Template  string    `json:"template" bson:"template"`
	Version   int       `json:"version" bson:"version"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	IsActive  bool      `json:"is_active" bson:"is_active"`
}
