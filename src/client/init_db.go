package main

import (
	"common"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func main() {
	common.InitConfig()
	config := common.Config

	quit := make(chan bool)

	dbURI := common.BuildURIFrom(&config.DB)
	db, _ := gorm.Open("postgres", dbURI)

	db.DB()
	db.LogMode(true)
	db.CreateTable(&(common.Article{}), &(common.CrawlLog{}))

	<-quit
}
