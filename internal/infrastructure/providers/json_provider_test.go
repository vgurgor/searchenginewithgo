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
			"items": []map[string]any{
				{
					"id": "vid_12345", "heading": "Sample Video Title", "type": "video",
					"summary": "Video description", "video_url": "https://example.com/video.mp4",
					"thumbnail": "https://example.com/thumb.jpg",
					"statistics": map[string]any{"view_count": 150000, "like_count": 5000},
					"release_date": "2024-11-01T10:00:00Z",
				},
				{
					"id": "art_67890", "heading": "Sample Article Title", "type": "article",
					"summary": "Article description", "article_url": "https://example.com/article",
					"read_duration": 8, "engagement": map[string]any{"reaction_count": 250},
					"release_date": "2024-10-15T14:30:00Z",
				},
			},
			"total": 2,
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
	if items[0].ContentType != "video" || items[1].ContentType != "text" {
		t.Fatalf("unexpected content types: %#v", items)
	}
}

func TestJSONProvider_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"items":[],"total":0}`))
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
		_, _ = w.Write([]byte(`{"items":[{"id":`)) // malformed
	}))
	defer srv.Close()
	p := NewJSONProvider(srv.URL, 2*time.Second)
	if _, err := p.FetchContents(); err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}


