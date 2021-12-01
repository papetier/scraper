package main

import (
	"github.com/papetier/scraper/pkg/config"
	"github.com/papetier/scraper/pkg/database"
	"github.com/papetier/scraper/pkg/scraper"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.LoadOrPrintVersion()

	// connect to the DB
	database.Connect()
	defer database.CloseConnection()

	// setup scrapers
	scraper.Setup()

	websiteList, err := database.GetWebsites()
	if err != nil {
		log.Fatalf("an error occurred fetching the website list: %v", err)
	}

	scraper.ScrapeWebsites(websiteList)
}
