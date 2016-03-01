package main

import (
	"common"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		helpMessage()
	}

	var tables = []interface{}{
		&common.Article{},
		&common.CrawlLog{}}

	switch os.Args[1] {
	case "create_all":
		db := configure()
		db.CreateTable(tables...)
	case "drop_all":
		db := configure()
		db.DropTable(tables...)
	case "reload_all":
		db := configure()
		db.DropTable(tables...)
		db.CreateTable(tables...)
	default:
		helpMessage()
	}
}

func helpMessage() {
	fmt.Println("./database [create_all, drop_all, reload_all]")
	os.Exit(1)
}

func configure() gorm.DB {
	common.InitConfig()
	config := common.Config
	dbURI := common.BuildURIFrom(&config.DB)
	db, _ := gorm.Open("postgres", dbURI)
	db.DB()
	db.LogMode(true)

	return db
}
