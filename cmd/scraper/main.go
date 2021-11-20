package main

import (
	"github.com/papetier/scraper/pkg/config"
	"github.com/papetier/scraper/pkg/database"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.Load()

	// connect to the DB
	database.Connect()
	defer database.CloseConnection()

	log.Infof("scraping...")
}
