package config

import (
	"github.com/spf13/viper"
	"time"
)

func setDefaultConfigValues() {
	// logger defaults
	viper.SetDefault("LOG_LEVEL", "info")

	// DB defaults
	viper.SetDefault("POSTGRES_DATABASE", "postgres")
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_PASSWORD", "postgres")
	viper.SetDefault("POSTGRES_USER", "postgres")

	// arXiv scraper defaults
	viper.SetDefault("ARXIV_ACCEPT_INSECURE_HTTP", false)
	viper.SetDefault("ARXIV_REQUEST_TIMEOUT", 30*time.Second)
	viper.SetDefault("ARXIV_DUPLICATED_THRESHOLD", 3)
	viper.SetDefault("ARXIV_MAX_RESULTS", 1000)
	viper.SetDefault("ARXIV_SEARCH_START", 0)
	viper.SetDefault("ARXIV_SORT_BY", "submittedDate")
	viper.SetDefault("ARXIV_SORT_ORDER", "ascending")
}
