package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/chat/completions") {
			http.Error(w, "bad path", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer k" {
			http.Error(w, "no auth", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"{\"ok\":true}"}}],"usage":{"total_tokens":42}}`))
	}))
	defer srv.Close()

	c := NewClient()
	resp, err := c.Chat(context.Background(), ChatRequest{
		BaseURL:     srv.URL,
		APIKey:      "k",
		Model:       "m",
		Temperature: 0.2,
		MaxTokens:   100,
		Messages:    []ChatMessage{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Content, "ok") {
		t.Fatalf("unexpected content: %s", resp.Content)
	}
	if resp.Usage != 42 {
		t.Fatalf("expected usage 42, got %d", resp.Usage)
	}
}

func TestChatError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("boom"))
	}))
	defer srv.Close()
	c := NewClient()
	if _, err := c.Chat(context.Background(), ChatRequest{BaseURL: srv.URL, APIKey: "k", Model: "m", Messages: []ChatMessage{{Role: "user", Content: "x"}}}); err == nil {
		t.Fatal("expected error")
	}
}
