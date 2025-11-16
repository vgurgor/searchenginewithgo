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
	Items []struct {
		ID         string `json:"id"`
		Heading    string `json:"heading"`
		Type       string `json:"type"`
		Summary    string `json:"summary"`
		VideoURL   string `json:"video_url,omitempty"`
		ArticleURL string `json:"article_url,omitempty"`
		Thumbnail  string `json:"thumbnail,omitempty"`
		ReadDur    int    `json:"read_duration,omitempty"`
		Stats      struct {
			ViewCount int64 `json:"view_count,omitempty"`
			LikeCount int64 `json:"like_count,omitempty"`
		} `json:"statistics,omitempty"`
		Engagement struct {
			ReactionCount int `json:"reaction_count,omitempty"`
		} `json:"engagement,omitempty"`
		ReleaseDate time.Time `json:"release_date"`
	} `json:"items"`
	Total int `json:"total"`
}

func (p *JSONProvider) FetchContents() ([]domainp.ProviderContent, error) {
	u, _ := url.Parse(p.BaseURL + "/contents")
	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", p.Limit))
	q.Set("offset", fmt.Sprintf("%d", p.Offset))
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d from provider1", resp.StatusCode)
	}
	var pr provider1Response
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	out := make([]domainp.ProviderContent, 0, len(pr.Items))
	for _, it := range pr.Items {
		pc := domainp.ProviderContent{
			ProviderID:        p.Provider,
			ProviderContentID: it.ID,
			Title:             it.Heading,
			Description:       it.Summary,
			PublishedAt:       it.ReleaseDate,
		}
		switch it.Type {
		case "video":
			pc.ContentType = "video"
			pc.URL = it.VideoURL
			pc.ThumbnailURL = it.Thumbnail
			if it.Stats.ViewCount != 0 {
				pc.Views = &it.Stats.ViewCount
			}
			if it.Stats.LikeCount != 0 {
				pc.Likes = &it.Stats.LikeCount
			}
		case "article":
			pc.ContentType = "text"
			pc.URL = it.ArticleURL
			if it.ReadDur != 0 {
				pc.ReadingTime = &it.ReadDur
			}
			if it.Engagement.ReactionCount != 0 {
				pc.Reactions = &it.Engagement.ReactionCount
			}
		default:
			pc.ContentType = "text"
		}
		out = append(out, pc)
	}
	return out, nil
}


