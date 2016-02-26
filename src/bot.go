package main

import (
	"encoding/json"
	"fmt"
	irc "github.com/fluffle/goirc/client"
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
	Id    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"url"`
}

func main() {
	duration := 3 * time.Second

	cfg := irc.NewConfig("kpop♡")
	cfg.Server = "irc.ozinger.org"
	cfg.NewNick = func(n string) string { return n + "♥" }
	c := irc.Client(cfg)
	ircChannel := "#freyja-test"
	feed := make([]Article, 30)

	quit := make(chan struct{})
	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) { conn.Join(ircChannel) })
	c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { quit <- struct{}{} })

	if err := c.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
	}

	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, article := range getKpopNews(&feed) {
					msg := "/r/kpop - [" + article.Title + "](" + article.Link + ")"
					c.Privmsg(ircChannel, msg)
					fmt.Println(msg)
				}
			}
		}
	}()

	<-quit
}

func getKpopNews(feed *[]Article) []Article {
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

	newArticles := getNewArticles(feed, &parsed)
	*feed = parsed
	return newArticles
}

func getNewArticles(existing *[]Article, parsed *[]Article) []Article {
	newArticles := *parsed

	for parsedIdx, parsedItem := range *parsed {
		for _, extItem := range *existing {
			if extItem == parsedItem {
				newArticles = append(newArticles[:parsedIdx], newArticles[parsedIdx+1:]...)
			}
		}
	}

	return newArticles
}
