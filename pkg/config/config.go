package config

import (
	"fmt"
	"github.com/papetier/scraper/pkg/logger"
	"github.com/papetier/scraper/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

type Environment string

const (
	DefaultEnvironment = "local"
	EnvironmentKey     = "SCRAPER_ENVIRONMENT"
)

func LoadOrPrintVersion() {
	// parse flags
	parseFlags()

	// print version and stop
	shouldPrintVersion := viper.GetBool("version")
	if shouldPrintVersion {
		fmt.Println(version.String())
		os.Exit(0)
	}

	// or continue loading the environment
	load()
}

func load() {
	// set environment
	environment := os.Getenv(EnvironmentKey)
	if environment == "" {
		environment = DefaultEnvironment
	}
	log.Debugf("config environment set to: %s", environment)

	// set default values
	setDefaultConfigValues()

	// load config from environment variables
	viper.AutomaticEnv()

	// define config file to read
	viper.SetConfigType("env")
	viper.SetConfigName(environment)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error when getting current working directory: %s", err)
	}
	viper.AddConfigPath(wd)

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			if environment != DefaultEnvironment {
				log.Fatalf("no config file found for environment: %s", environment)
			}
		} else {
			// Config file was found but another error occurred
			log.Fatal(err)
		}
	}

	// logger
	logger.Configure(viper.GetString("LOG_LEVEL"))

	// DB config
	loadDbConfig()

	log.Infof("%s config successfully loaded", environment)
}

func parseFlags() {
	pflag.BoolP("version", "v", false, "prints the version")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("could not bind the flags to the config: %s", err)
	}
}
