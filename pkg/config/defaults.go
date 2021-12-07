package config

import "github.com/spf13/viper"

func setDefaultConfigValues() {
	// logger defaults
	viper.SetDefault("LOG_LEVEL", "info")

	// DB defaults
	viper.SetDefault("POSTGRES_DATABASE", "postgres")
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_PASSWORD", "postgres")
	viper.SetDefault("POSTGRES_USER", "postgres")

	// scraper defaults
	viper.SetDefault("SCRAPER_ACCEPT_INSECURE_HTTP", false)
}
