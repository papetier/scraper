package config

import (
	"github.com/spf13/viper"
	"time"
)

type ScraperConfig struct {
	IsInsecureHttpAccepted bool
	RequestTimeout         time.Duration
}

var Scraper *ScraperConfig

func loadScraperConfig() {
	Scraper = &ScraperConfig{
		IsInsecureHttpAccepted: viper.GetBool("SCRAPER_ACCEPT_INSECURE_HTTP"),
		RequestTimeout:         viper.GetDuration("SCRAPER_REQUEST_TIMEOUT"),
	}
}
