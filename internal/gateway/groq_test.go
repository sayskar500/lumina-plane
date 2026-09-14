package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newMockGroq starts a test server that mimics the Groq chat completions API.
func newMockGroq(t *testing.T, status int, modelOut *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status == http.StatusOK {
			var body struct {
				Model string `json:"model"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			*modelOut = body.Model

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"choices": []map[string]interface{}{
					{"message": map[string]string{"content": "mock answer"}},
				},
				"usage": map[string]int{"total_tokens": 42},
			})
			return
		}
		http.Error(w, "error", status)
	}))
}

func TestGenerateUsesRequestedModel(t *testing.T) {
	var gotModel string
	srv := newMockGroq(t, http.StatusOK, &gotModel)
	defer srv.Close()

	p := NewGroqProvider("test-key", "llama3-8b-8192")
	p.BaseURL = srv.URL

	resp, err := p.Generate(context.Background(), "hi", "llama3-70b-8192")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if gotModel != "llama3-70b-8192" {
		t.Errorf("request model = %q, want llama3-70b-8192", gotModel)
	}
	if resp.ModelUsed != "llama3-70b-8192" {
		t.Errorf("ModelUsed = %q, want llama3-70b-8192", resp.ModelUsed)
	}
	if resp.TokensUsed != 42 {
		t.Errorf("TokensUsed = %d, want 42", resp.TokensUsed)
	}
	if resp.Text != "mock answer" {
		t.Errorf("Text = %q, want mock answer", resp.Text)
	}
}

func TestGenerateFallsBackToDefaultModel(t *testing.T) {
	var gotModel string
	srv := newMockGroq(t, http.StatusOK, &gotModel)
	defer srv.Close()

	p := NewGroqProvider("test-key", "llama3-8b-8192")
	p.BaseURL = srv.URL

	resp, err := p.Generate(context.Background(), "hi", "")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if gotModel != "llama3-8b-8192" {
		t.Errorf("request model = %q, want default llama3-8b-8192", gotModel)
	}
	if resp.ModelUsed != "llama3-8b-8192" {
		t.Errorf("ModelUsed = %q, want llama3-8b-8192", resp.ModelUsed)
	}
}

func TestGenerateMapsAuthError(t *testing.T) {
	srv := newMockGroq(t, http.StatusUnauthorized, nil)
	defer srv.Close()

	p := NewGroqProvider("bad-key", "llama3-8b-8192")
	p.BaseURL = srv.URL

	_, err := p.Generate(context.Background(), "hi", "")
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("err = %v, want a 401 Unauthorized error", err)
	}
}

func TestGenerateFailsWithoutAPIKey(t *testing.T) {
	p := NewGroqProvider("", "llama3-8b-8192")
	if _, err := p.Generate(context.Background(), "hi", ""); err == nil {
		t.Fatal("expected an error when the API key is missing")
	}
}
