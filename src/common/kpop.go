package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Feed struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash  string `json:"modhash"`
		Children []struct {
			Kind string  `json:"kind"`
			Data Article `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Article struct {
	ID    string `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
	Link  string `json:"url"`
}

type CrawlLog struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `sql:"default:now()"`
}

func GetKpopNews() []Article {
	resp, feedErr := http.Get("https://www.reddit.com/r/kpop/hot.json")
	if feedErr != nil {
		fmt.Printf("Error: %s\n", feedErr.Error())
	}
	defer resp.Body.Close()

	newFeed := Feed{}
	if jsonErr := json.NewDecoder(resp.Body).Decode(&newFeed); jsonErr != nil {
		fmt.Printf("Error: %s\n", jsonErr.Error())
	}

	var parsed []Article
	for _, children := range newFeed.Data.Children {
		parsed = append(parsed, children.Data)
	}

	return parsed
}
