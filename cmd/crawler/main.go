package main

import (
	"github.com/papetier/crawler/pkg/config"
	"github.com/papetier/crawler/pkg/database"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.Load()

	// connect to the DB
	database.Connect()
	defer database.CloseConnection()

	log.Infof("crawling...")
}
