package handlers

import (
	"encoding/xml"
	"net/http"
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

func MockProvider1Handler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	all := make([]provider1Item, 0, 40)
	now := time.Now().UTC()
	for i := 0; i < 40; i++ {
		if i%2 == 0 {
			// video
			it := provider1Item{
				ID:         "vid_" + strconv.Itoa(10000+i),
				Heading:    "Sample Video Title " + strconv.Itoa(i),
				Type:       "video",
				Summary:    "Video description",
				VideoURL:   "https://example.com/video" + strconv.Itoa(i) + ".mp4",
				Thumbnail:  "https://example.com/thumb" + strconv.Itoa(i) + ".jpg",
				ReleaseDate: now.Add(-time.Duration(i) * 24 * time.Hour),
			}
			it.Stats.ViewCount = int64(1000 + i*73)
			it.Stats.LikeCount = int64(100 + i*5)
			all = append(all, it)
		} else {
			// article
			it := provider1Item{
				ID:         "art_" + strconv.Itoa(10000+i),
				Heading:    "Sample Article Title " + strconv.Itoa(i),
				Type:       "article",
				Summary:    "Article description",
				ArticleURL: "https://example.com/article/" + strconv.Itoa(i),
				ReadDur:    5 + (i % 20),
				ReleaseDate: now.Add(-time.Duration(i) * 24 * time.Hour),
			}
			it.Engagement.ReactionCount = 50 + (i % 200)
			all = append(all, it)
		}
	}
	end := offset + limit
	if offset > len(all) {
		offset = len(all)
	}
	if end > len(all) {
		end = len(all)
	}
	c.JSON(http.StatusOK, provider1Resp{
		Items: all[offset:end],
		Total: len(all),
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

func MockProvider2Handler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 50 {
		size = 50
	}
	now := time.Now().UTC()
	all := make([]Provider2Item, 0, 40)
	for i := 0; i < 40; i++ {
		if i%2 == 0 {
			v := int64(5000 + i*77)
			l := int64(200 + i*3)
			all = append(all, Provider2Item{
				ContentID: "v_" + strconv.Itoa(9000+i),
				Title:     "Another Video Title " + strconv.Itoa(i),
				Category:  "video",
				Desc:      "Video description here",
				Link:      "https://example.com/video" + strconv.Itoa(i) + ".mp4",
				Image:     "https://example.com/thumb" + strconv.Itoa(i) + ".jpg",
				Metrics:   Provider2Met{Views: &v, Likes: &l},
				PubDate:   now.Add(-time.Duration(i) * 24 * time.Hour).Format(time.RFC3339),
			})
		} else {
			r := 100 + (i % 300)
			all = append(all, Provider2Item{
				ContentID: "a_" + strconv.Itoa(9000+i),
				Title:     "Blog Post Title " + strconv.Itoa(i),
				Category:  "text",
				Desc:      "Blog post description",
				Link:      "https://example.com/blog/" + strconv.Itoa(i),
				ReadTime:  5 + (i % 20),
				Metrics:   Provider2Met{Reactions: &r},
				PubDate:   now.Add(-time.Duration(i) * 24 * time.Hour).Format(time.RFC3339),
			})
		}
	}
	start := (page - 1) * size
	end := start + size
	if start > len(all) {
		start = len(all)
	}
	if end > len(all) {
		end = len(all)
	}
	feed := Provider2Feed{Items: all[start:end]}
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.XML(http.StatusOK, feed)
}


