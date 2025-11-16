package providers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestXMLProvider_FetchContents(t *testing.T) {
	xmlBody := `<?xml version="1.0" encoding="UTF-8"?>
<feed>
  <items>
    <item>
      <id>v1</id>
      <headline>Introduction to Docker</headline>
      <type>video</type>
      <stats>
        <views>22000</views>
        <likes>1800</likes>
        <duration>25:15</duration>
      </stats>
      <publication_date>2024-03-15</publication_date>
    </item>
    <item>
      <id>a1</id>
      <headline>Clean Architecture in Go</headline>
      <type>article</type>
      <stats>
        <reading_time>8</reading_time>
        <reactions>450</reactions>
        <comments>25</comments>
      </stats>
      <publication_date>2024-03-14</publication_date>
    </item>
  </items>
  <meta>
    <total_count>75</total_count>
    <current_page>1</current_page>
    <items_per_page>10</items_per_page>
  </meta>
</feed>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(xmlBody))
	}))
	defer srv.Close()

	p := NewXMLProvider(srv.URL, 5*time.Second)
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
	if items[0].Title != "Introduction to Docker" {
		t.Fatalf("unexpected title: %s", items[0].Title)
	}
}

func TestXMLProvider_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte(`<feed></feed>`))
	}))
	defer srv.Close()
	p := NewXMLProvider(srv.URL, 50*time.Millisecond)
	if _, err := p.FetchContents(); err == nil {
		t.Fatalf("expected timeout")
	}
}
