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
  <content>
    <content_id>v_9876</content_id>
    <title>Another Video Title</title>
    <category>video</category>
    <desc>Video description here</desc>
    <link>https://example.com/video2.mp4</link>
    <image>https://example.com/thumb2.jpg</image>
    <metrics><views>85000</views><likes>3200</likes></metrics>
    <pub_date>2024-11-10T08:00:00Z</pub_date>
  </content>
  <content>
    <content_id>a_5432</content_id>
    <title>Blog Post Title</title>
    <category>text</category>
    <desc>Blog post description</desc>
    <link>https://example.com/blog</link>
    <read_time>12</read_time>
    <metrics><reactions>180</reactions></metrics>
    <pub_date>2024-09-20T16:45:00Z</pub_date>
  </content>
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
	if items[0].ContentType != "video" || items[1].ContentType != "text" {
		t.Fatalf("unexpected content types: %#v", items)
	}
}


