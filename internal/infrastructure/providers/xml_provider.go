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
	XMLName xml.Name   `xml:"feed"`
	Items   []xmlItem  `xml:"content"`
}
type xmlItem struct {
	ContentID string    `xml:"content_id"`
	Title     string    `xml:"title"`
	Category  string    `xml:"category"`
	Desc      string    `xml:"desc"`
	Link      string    `xml:"link"`
	Image     string    `xml:"image"`
	ReadTime  int       `xml:"read_time"`
	Metrics   xmlMetric `xml:"metrics"`
	PubDate   string    `xml:"pub_date"`
}
type xmlMetric struct {
	Views     *int64 `xml:"views"`
	Likes     *int64 `xml:"likes"`
	Reactions *int   `xml:"reactions"`
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
			ProviderContentID: it.ContentID,
			Title:             it.Title,
			Description:       it.Desc,
			URL:               it.Link,
			ThumbnailURL:      it.Image,
		}
		// parse time
		if ts, err := time.Parse(time.RFC3339, it.PubDate); err == nil {
			pc.PublishedAt = ts
		} else {
			pc.PublishedAt = time.Now().UTC()
		}
		switch it.Category {
		case "video":
			pc.ContentType = "video"
			pc.Views = it.Metrics.Views
			pc.Likes = it.Metrics.Likes
		default:
			pc.ContentType = "text"
			if it.ReadTime != 0 {
				pc.ReadingTime = &it.ReadTime
			}
			pc.Reactions = it.Metrics.Reactions
		}
		out = append(out, pc)
	}
	return out, nil
}


