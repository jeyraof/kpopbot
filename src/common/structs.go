package common

import (
	"fmt"
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
	return fmt.Sprintf("%s - [%s](%s) #kpop #k_pop", target, article.Title, article.Link)
}
