package handlers

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Provider1 JSON mock
type provider1Item struct {
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
}

type provider1Resp struct {
	Items []provider1Item `json:"items"`
	Total int             `json:"total"`
}

// File-backed schema for provider1.json at repo root
type provider1File struct {
	Contents []struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Type    string `json:"type"`
		Metrics struct {
			Views       int64  `json:"views"`
			Likes       int64  `json:"likes"`
			Duration    string `json:"duration,omitempty"`
			ReadingTime int    `json:"reading_time,omitempty"`
			Reactions   int    `json:"reactions,omitempty"`
		} `json:"metrics"`
		PublishedAt time.Time `json:"published_at"`
	} `json:"contents"`
	Pagination struct {
		Total   int `json:"total"`
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	} `json:"pagination"`
}

func MockProvider1Handler(c *gin.Context) {
	// File-only mode: require file, otherwise fail
	fileOnly := os.Getenv("PROVIDERS_FILE_ONLY") == "true"

	// Try to serve from file if path is provided
	if fp := os.Getenv("PROVIDER1_FILE_PATH"); fp != "" {
		if f, err := os.Open(fp); err == nil {
			defer f.Close()
			var raw provider1File
			if err := json.NewDecoder(f).Decode(&raw); err == nil {
				// Return the file content as-is (in the new format)
				c.JSON(http.StatusOK, raw)
				return
			}
		}
	}

	if fileOnly {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": "FILE_REQUIRED", "message": "PROVIDER1_FILE_PATH tanımlı değil veya dosya okunamadı"},
		})
		return
	}

	// Fallback: Generate synthetic data in the new format
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	now := time.Now().UTC()
	contents := make([]struct {
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
	}, 0, limit)

	for i := 0; i < limit && i < 40; i++ {
		item := struct {
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
		}{}

		if i%2 == 0 {
			item.ID = "v" + strconv.Itoa(100+i)
			item.Title = "Sample Video " + strconv.Itoa(i)
			item.Type = "video"
			item.Metrics.Views = int64(10000 + i*100)
			item.Metrics.Likes = int64(800 + i*10)
			item.Metrics.Duration = "15:30"
			item.Tags = []string{"programming", "tutorial"}
		} else {
			item.ID = "a" + strconv.Itoa(100+i)
			item.Title = "Sample Article " + strconv.Itoa(i)
			item.Type = "article"
			item.Metrics.Reactions = 100 + i*5
			item.Tags = []string{"programming", "article"}
		}
		item.PublishedAt = now.Add(-time.Duration(i) * 24 * time.Hour)
		contents = append(contents, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
		"pagination": gin.H{
			"total":    150,
			"page":     1,
			"per_page": limit,
		},
	})
}

// Provider2 XML mock
type Provider2Feed struct {
	XMLName xml.Name        `xml:"feed"`
	Items   []Provider2Item `xml:"content"`
}
type Provider2Item struct {
	ContentID string       `xml:"content_id"`
	Title     string       `xml:"title"`
	Category  string       `xml:"category"`
	Desc      string       `xml:"desc"`
	Link      string       `xml:"link"`
	Image     string       `xml:"image,omitempty"`
	ReadTime  int          `xml:"read_time,omitempty"`
	Metrics   Provider2Met `xml:"metrics"`
	PubDate   string       `xml:"pub_date"`
}
type Provider2Met struct {
	Views     *int64 `xml:"views,omitempty"`
	Likes     *int64 `xml:"likes,omitempty"`
	Reactions *int   `xml:"reactions,omitempty"`
}

// File-backed schema for provider2.xml at repo root
type provider2File struct {
	XMLName xml.Name `xml:"feed"`
	Items   struct {
		Item []struct {
			ID       string `xml:"id"`
			Headline string `xml:"headline"`
			Type     string `xml:"type"`
			Stats    struct {
				Views       *int64 `xml:"views"`
				Likes       *int64 `xml:"likes"`
				Duration    string `xml:"duration"`
				ReadingTime *int   `xml:"reading_time"`
				Reactions   *int   `xml:"reactions"`
				Comments    *int   `xml:"comments"`
			} `xml:"stats"`
			PublicationDate string `xml:"publication_date"`
		} `xml:"item"`
	} `xml:"items"`
	Meta struct {
		TotalCount   int `xml:"total_count"`
		CurrentPage  int `xml:"current_page"`
		ItemsPerPage int `xml:"items_per_page"`
	} `xml:"meta"`
}

func MockProvider2Handler(c *gin.Context) {
	// File-only mode: require file, otherwise fail
	fileOnly := os.Getenv("PROVIDERS_FILE_ONLY") == "true"

	// Try to serve from file if path is provided
	if fp := os.Getenv("PROVIDER2_FILE_PATH"); fp != "" {
		if f, err := os.Open(fp); err == nil {
			defer f.Close()
			var raw provider2File
			if err := xml.NewDecoder(f).Decode(&raw); err == nil {
				// Return the file content as-is (in the new format)
				c.Header("Content-Type", "application/xml; charset=utf-8")
				c.XML(http.StatusOK, raw)
				return
			}
		}
	}

	if fileOnly {
		c.Header("Content-Type", "application/xml; charset=utf-8")
		c.String(http.StatusInternalServerError, `<error code="FILE_REQUIRED">PROVIDER2_FILE_PATH tanımlı değil veya dosya okunamadı</error>`)
		return
	}

	// Fallback: Generate synthetic data in the new format
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if size <= 0 {
		size = 20
	}
	if size > 50 {
		size = 50
	}

	now := time.Now().UTC()
	type xmlItem struct {
		ID       string `xml:"id"`
		Headline string `xml:"headline"`
		Type     string `xml:"type"`
		Stats    struct {
			Views       *int64 `xml:"views,omitempty"`
			Likes       *int64 `xml:"likes,omitempty"`
			Duration    string `xml:"duration,omitempty"`
			ReadingTime *int   `xml:"reading_time,omitempty"`
			Reactions   *int   `xml:"reactions,omitempty"`
		} `xml:"stats"`
		PublicationDate string `xml:"publication_date"`
	}

	items := make([]xmlItem, 0, size)
	for i := 0; i < size && i < 20; i++ {
		item := xmlItem{}
		if i%2 == 0 {
			item.ID = "v" + strconv.Itoa(200+i)
			item.Headline = "Sample Video " + strconv.Itoa(i)
			item.Type = "video"
			v := int64(15000 + i*200)
			l := int64(1200 + i*15)
			item.Stats.Views = &v
			item.Stats.Likes = &l
			item.Stats.Duration = "20:30"
		} else {
			item.ID = "a" + strconv.Itoa(200+i)
			item.Headline = "Sample Article " + strconv.Itoa(i)
			item.Type = "article"
			rt := 8
			r := 300 + i*10
			item.Stats.ReadingTime = &rt
			item.Stats.Reactions = &r
		}
		item.PublicationDate = now.Add(-time.Duration(i) * 24 * time.Hour).Format("2006-01-02")
		items = append(items, item)
	}

	type feed struct {
		XMLName xml.Name `xml:"feed"`
		Items   struct {
			Item []xmlItem `xml:"item"`
		} `xml:"items"`
		Meta struct {
			TotalCount   int `xml:"total_count"`
			CurrentPage  int `xml:"current_page"`
			ItemsPerPage int `xml:"items_per_page"`
		} `xml:"meta"`
	}

	response := feed{}
	response.Items.Item = items
	response.Meta.TotalCount = 75
	response.Meta.CurrentPage = 1
	response.Meta.ItemsPerPage = size

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.XML(http.StatusOK, response)
}
