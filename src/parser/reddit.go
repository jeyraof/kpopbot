package parser

import (
	"common"
	"encoding/json"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
	"net/http"
	"time"
)

type RedditFeed struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash  string `json:"modhash"`
		Children []struct {
			Kind string         `json:"kind"`
			Data common.Article `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func GetRedditNews(target string) []common.Article {
	client := &http.Client{}

	url := fmt.Sprintf("https://www.reddit.com%s/hot.json", target)
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		fmt.Printf("Error: %s\n", reqErr.Error())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36")
	resp, feedErr := client.Do(req)
	if feedErr != nil {
		fmt.Printf("Error: %s\n", feedErr.Error())
	}
	defer resp.Body.Close()

	newFeed := RedditFeed{}
	if jsonErr := json.NewDecoder(resp.Body).Decode(&newFeed); jsonErr != nil {
		fmt.Printf("Error: %s\n", jsonErr.Error())
	}

	var parsed []common.Article
	for _, children := range newFeed.Data.Children {
		parsed = append(parsed, children.Data)
	}

	return parsed
}

func RedditRoutine(
	config *common.ConfigType,
	db *gorm.DB, irc *irc.Conn,
	crawlerQuit <-chan struct{},
	period uint,
	target string) {

	ticker := time.NewTicker(time.Duration(period) * time.Second)
	for {
		select {
		case <-ticker.C:
			now := common.CrawlLog{Target: target, CreatedAt: time.Now()}
			db.Create(&now)

			articles := GetRedditNews(target)
			for _, article := range articles {
				article.Target = target
				common.ArticleShorten(&config.Google, &article)
                common.ArticleUnescape(&article)
				if err := db.Create(&article).Error; err != nil {
					// Error Code Reference: https://github.com/lib/pq/blob/master/error.go#L78
					switch err.(*pq.Error).Code.Name() {
					case "unique_violation":
						// TODO: Handle integrity error on unique constraint
					}
				} else {
					msg := common.BuildMessage(target, &article)
					irc.Privmsg(config.IRC.Channel, msg)
					common.UpdateStatus(&config.Twitter, msg)
				}

			}
		case <-crawlerQuit:
			ticker.Stop()
			return
		}
	}
}
