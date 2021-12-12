package config

import (
	"github.com/spf13/viper"
	"time"
)

type ArxivConfig struct {
	IsInsecureHttpAccepted bool
	RequestTimeout         time.Duration
	DuplicatedThreshold    int
	MaxResults             int
	SearchStart            int
	SortBy                 string
	SortOrder              string
}

var Arxiv *ArxivConfig

func loadScraperConfig() {
	// arXiv config
	Arxiv = &ArxivConfig{
		IsInsecureHttpAccepted: viper.GetBool("ARXIV_ACCEPT_INSECURE_HTTP"),
		RequestTimeout:         viper.GetDuration("ARXIV_REQUEST_TIMEOUT"),
		DuplicatedThreshold:    viper.GetInt("ARXIV_DUPLICATED_THRESHOLD"),
		MaxResults:             viper.GetInt("ARXIV_MAX_RESULTS"),
		SearchStart:            viper.GetInt("ARXIV_SEARCH_START"),
		SortBy:                 viper.GetString("ARXIV_SORT_BY"),
		SortOrder:              viper.GetString("ARXIV_SORT_ORDER"),
	}
}
