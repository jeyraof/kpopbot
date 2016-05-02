package common

import (
	"fmt"
	"html"
	"time"
)

type Article struct {
	ID        string    `gorm:"primary_key" json:"id" xml:"guid"`
	Target    string    `gorm:"primary_key"`
	Title     string    `json:"title" xml:"title"`
	Link      string    `json:"url" xml:"link"`
	CrawledAt time.Time `sql:"default:now()"`
}

type CrawlLog struct {
	ID        uint `gorm:"primary_key"`
	Target    string
	CreatedAt time.Time `sql:"default:now()"`
}

func BuildMessage(target string, article *Article) string {
	return fmt.Sprintf("%s %s — %s", article.Title, article.Link, target)
}

func ArticleUnescape(article *Article) {
	article.Title = html.UnescapeString(article.Title)
	article.Link = html.UnescapeString(article.Link)
}
