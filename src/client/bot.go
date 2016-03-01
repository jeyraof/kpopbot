package main

import (
	"common"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
	"time"
)

func main() {
	common.InitConfig()
	config := common.Config

	// Configure Database
	dbURI := common.BuildURIFrom(&config.DB)
	db, dbErr := gorm.Open("postgres", dbURI)
	if dbErr != nil {
		fmt.Printf("Database error: %s\n", dbErr.Error())
	}

	// Configure IRC
	cfg := irc.NewConfig(config.IRC.Nick)
	cfg.Server = config.IRC.Server
	cfg.NewNick = func(n string) string { return config.IRC.NickAlternative }
	c := irc.Client(cfg)
	ircQuit := make(chan struct{})

	// Handler for IRC actions
	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) { conn.Join(config.IRC.Channel) })
	c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { ircQuit <- struct{}{} })

	// IRC Connect!
	if ircErr := c.Connect(); ircErr != nil {
		fmt.Printf("Connection error: %s\n", ircErr.Error())
	}

	// Periodic Crawl Kpopnews
	crawlerQuit := make(chan struct{})
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				now := common.CrawlLog{CreatedAt: time.Now()}
				db.Create(&now)

				articles := common.GetKpopNews()
				for _, article := range articles {
					common.ArticleShorten(&config.Google, &article)
					if err := db.Create(&article).Error; err != nil {
						switch err.(*pq.Error).Code.Name() {
						case "unique_violation":
							// TODO: Handle integrity error on unique constraint
						}
					} else {
						SendMsg(&config.IRC, c, &article)
					}

				}
			case <-crawlerQuit:
				ticker.Stop()
				return
			}
		}
	}()

	<-ircQuit
}

func SendMsg(config *common.IRCConfigType, irc *irc.Conn, article *common.Article) {
	msg := "/r/kpop - [" + (*article).Title + "](" + (*article).Link + ")"
	(*irc).Privmsg((*config).Channel, msg)
}
