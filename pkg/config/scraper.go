package config

import "github.com/spf13/viper"

type ScraperConfig struct {
	IsInsecureHttpAccepted bool
}

var Scraper *ScraperConfig

func loadScraperConfig() {
	Scraper = &ScraperConfig{
		IsInsecureHttpAccepted: viper.GetBool("SCRAPER_ACCEPT_INSECURE_HTTP"),
	}
}
