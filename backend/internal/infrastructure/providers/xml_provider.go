package providers

import (
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	domainp "search_engine/internal/domain/providers"
)

type XMLProvider struct {
	Client   *http.Client
	BaseURL  string
	Provider string
	Page     int
	Size     int
}

func NewXMLProvider(baseURL string, timeout time.Duration) *XMLProvider {
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
	return &XMLProvider{Client: client, BaseURL: baseURL, Provider: "provider2", Page: 1, Size: 40}
}

func (p *XMLProvider) GetProviderID() string { return p.Provider }
func (p *XMLProvider) GetRateLimit() domainp.RateLimit {
	return domainp.RateLimit{RequestsPerMinute: 80}
}

type xmlFeed struct {
	XMLName xml.Name  `xml:"feed"`
	Items   []xmlItem `xml:"items>item"`
}
type xmlItem struct {
	ID       string   `xml:"id"`
	Headline string   `xml:"headline"`
	Type     string   `xml:"type"`
	Stats    xmlStats `xml:"stats"`
	PubDate  string   `xml:"publication_date"`
}
type xmlStats struct {
	Views       *int64 `xml:"views"`
	Likes       *int64 `xml:"likes"`
	Duration    string `xml:"duration"`
	ReadingTime *int   `xml:"reading_time"`
	Reactions   *int   `xml:"reactions"`
	Comments    *int   `xml:"comments"`
}

func (p *XMLProvider) FetchContents() ([]domainp.ProviderContent, error) {
	u, _ := url.Parse(p.BaseURL + "/feed")
	q := u.Query()
	q.Set("page", fmt.Sprintf("%d", p.Page))
	q.Set("size", fmt.Sprintf("%d", p.Size))
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d from provider2", resp.StatusCode)
	}
	var feed xmlFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}
	out := make([]domainp.ProviderContent, 0, len(feed.Items))
	for _, it := range feed.Items {
		pc := domainp.ProviderContent{
			ProviderID:        p.Provider,
			ProviderContentID: it.ID,
			Title:             it.Headline,
			Description:       "", // No description in this provider format
		}
		// parse time - try multiple formats
		var ts time.Time
		if t, err := time.Parse(time.RFC3339, it.PubDate); err == nil {
			ts = t
		} else if t, err := time.Parse("2006-01-02", it.PubDate); err == nil {
			ts = t
		} else {
			ts = time.Now().UTC()
		}
		pc.PublishedAt = ts

		switch it.Type {
		case "video":
			pc.ContentType = "video"
			pc.URL = fmt.Sprintf("https://example.com/video/%s", it.ID)
			pc.ThumbnailURL = fmt.Sprintf("https://example.com/thumb/%s.jpg", it.ID)
			pc.Views = it.Stats.Views
			pc.Likes = it.Stats.Likes
		case "article":
			pc.ContentType = "text"
			pc.URL = fmt.Sprintf("https://example.com/article/%s", it.ID)
			pc.ReadingTime = it.Stats.ReadingTime
			pc.Reactions = it.Stats.Reactions
		default:
			pc.ContentType = "text"
			pc.URL = fmt.Sprintf("https://example.com/content/%s", it.ID)
		}
		out = append(out, pc)
	}
	return out, nil
}
