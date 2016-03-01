package main

import (
	"common"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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

	<-ircQuit
}
