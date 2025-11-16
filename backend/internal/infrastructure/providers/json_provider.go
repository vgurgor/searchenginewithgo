package providers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	domainp "search_engine/internal/domain/providers"
)

type JSONProvider struct {
	Client   *http.Client
	BaseURL  string
	Provider string
	Limit    int
	Offset   int
}

func NewJSONProvider(baseURL string, timeout time.Duration) *JSONProvider {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
		},
	}
	return &JSONProvider{Client: client, BaseURL: baseURL, Provider: "provider1", Limit: 40, Offset: 0}
}

func (p *JSONProvider) GetProviderID() string { return p.Provider }
func (p *JSONProvider) GetRateLimit() domainp.RateLimit {
	return domainp.RateLimit{RequestsPerMinute: 100}
}

type provider1Response struct {
	Contents []struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Type    string `json:"type"`
		Metrics struct {
			Views     int64  `json:"views,omitempty"`
			Likes     int64  `json:"likes,omitempty"`
			Duration  string `json:"duration,omitempty"`
			Reactions int    `json:"reactions,omitempty"`
		} `json:"metrics"`
		PublishedAt time.Time `json:"published_at"`
		Tags        []string  `json:"tags,omitempty"`
	} `json:"contents"`
	Pagination struct {
		Total   int `json:"total"`
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	} `json:"pagination,omitempty"`
}

func (p *JSONProvider) FetchContents() ([]domainp.ProviderContent, error) {
	u, _ := url.Parse(p.BaseURL + "/contents")
	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", p.Limit))
	q.Set("offset", fmt.Sprintf("%d", p.Offset))
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), http.NoBody)
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d from provider1", resp.StatusCode)
	}
	var pr provider1Response
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	out := make([]domainp.ProviderContent, 0, len(pr.Contents))
	for _, it := range pr.Contents {
		pc := domainp.ProviderContent{
			ProviderID:        p.Provider,
			ProviderContentID: it.ID,
			Title:             it.Title,
			Description:       "", // No description in this provider format
			PublishedAt:       it.PublishedAt,
		}
		switch it.Type {
		case "video":
			pc.ContentType = "video"
			pc.URL = fmt.Sprintf("https://example.com/video/%s", it.ID)
			pc.ThumbnailURL = fmt.Sprintf("https://example.com/thumb/%s.jpg", it.ID)
			if it.Metrics.Views != 0 {
				pc.Views = &it.Metrics.Views
			}
			if it.Metrics.Likes != 0 {
				pc.Likes = &it.Metrics.Likes
			}
		case "article":
			pc.ContentType = "text"
			pc.URL = fmt.Sprintf("https://example.com/article/%s", it.ID)
			if it.Metrics.Reactions != 0 {
				pc.Reactions = &it.Metrics.Reactions
			}
			// Estimate reading time from duration if available (not present in current format)
			readTime := 5 // default
			pc.ReadingTime = &readTime
		default:
			pc.ContentType = "text"
			pc.URL = fmt.Sprintf("https://example.com/content/%s", it.ID)
		}
		out = append(out, pc)
	}
	return out, nil
}
