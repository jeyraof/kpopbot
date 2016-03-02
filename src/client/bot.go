package main

import (
	"common"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/jinzhu/gorm"
	"parser"
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
	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(config.IRC.Channel)
		if config.IRC.Nickserv != "" {
			conn.Privmsg(config.IRC.Nickserv, fmt.Sprintf("로그인 %s %s", config.IRC.NickservNick, config.IRC.NickservPassword))
		}
	})
	c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { ircQuit <- struct{}{} })
	c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		message := line.Args[1]
		me := conn.Me()

		if len(message) >= len(me.Nick) && message[:len(me.Nick)] == me.Nick {
			conn.Privmsg(config.IRC.Channel, config.Repository)
		}
	})

	// IRC Connect!
	if ircErr := c.Connect(); ircErr != nil {
		fmt.Printf("Connection error: %s\n", ircErr.Error())
	}

	// Periodic Crawl Kpopnews
	crawlerQuit := make(chan struct{})
	go parser.RedditRoutine(&config, &db, c, crawlerQuit, 300, "/r/kpop")
	go parser.IdologyRoutine(&config, &db, c, crawlerQuit, 3600)

	<-ircQuit
}
