package parser

import (
	"common"
	"encoding/xml"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
	"net/http"
	"time"
)

type IdologyFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		XMLName  xml.Name         `xml:"channel"`
		ItemList []common.Article `xml:"item"`
	} `xml:"channel"`
}

func GetIdologyNews() []common.Article {
	resp, feedErr := http.Get("http://idology.kr/feed")
	if feedErr != nil {
		fmt.Printf("Error: %s\n", feedErr.Error())
	}
	defer resp.Body.Close()

	newFeed := IdologyFeed{}
	if xmlErr := xml.NewDecoder(resp.Body).Decode(&newFeed); xmlErr != nil {
		fmt.Printf("Error: %s\n", xmlErr.Error())
	}

	var parsed []common.Article
	for _, children := range newFeed.Channel.ItemList {
		parsed = append(parsed, children)
	}

	return parsed
}

func IdologyRoutine(
	config *common.ConfigType,
	db *gorm.DB, irc *irc.Conn,
	crawlerQuit <-chan struct{},
	period uint) {

	target := "/idology"
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	for {
		select {
		case <-ticker.C:
			now := common.CrawlLog{Target: target, CreatedAt: time.Now()}
			db.Create(&now)

			articles := GetIdologyNews()
			for _, article := range articles {
				article.Target = target
				common.ArticleShorten(&config.Google, &article)
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
