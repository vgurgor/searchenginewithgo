package providers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJSONProvider_FetchContents(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{
			"contents": []map[string]any{
				{
					"id": "v1", "title": "Go Programming Tutorial", "type": "video",
					"metrics":      map[string]any{"views": 15000, "likes": 1200, "duration": "15:30"},
					"published_at": "2024-03-15T10:00:00Z",
					"tags":         []string{"programming", "tutorial"},
				},
				{
					"id": "a1", "title": "Clean Code Article", "type": "article",
					"metrics":      map[string]any{"reactions": 450},
					"published_at": "2024-03-14T14:30:00Z",
					"tags":         []string{"programming", "article"},
				},
			},
			"pagination": map[string]any{
				"total": 150, "page": 1, "per_page": 10,
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	p := NewJSONProvider(srv.URL, 5*time.Second)
	items, err := p.FetchContents()
	if err != nil {
		t.Fatalf("FetchContents error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items got %d", len(items))
	}
	if items[0].ContentType != "video" {
		t.Fatalf("expected first item to be video, got %s", items[0].ContentType)
	}
	if items[1].ContentType != "text" {
		t.Fatalf("expected second item to be text, got %s", items[1].ContentType)
	}
	if items[0].Title != "Go Programming Tutorial" {
		t.Fatalf("unexpected title: %s", items[0].Title)
	}
}

func TestJSONProvider_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"contents":[],"pagination":{"total":0,"page":1,"per_page":10}}`))
	}))
	defer srv.Close()
	p := NewJSONProvider(srv.URL, 50*time.Millisecond)
	if _, err := p.FetchContents(); err == nil {
		t.Fatalf("expected timeout error")
	}
}

func TestJSONProvider_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"contents":[{"id":`)) // malformed
	}))
	defer srv.Close()
	p := NewJSONProvider(srv.URL, 2*time.Second)
	if _, err := p.FetchContents(); err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}
