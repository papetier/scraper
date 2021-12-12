package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type ArxivConfig struct {
	CategoryList           []string
	InitUrlList            []string
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
	// arXiv's category list
	var categoryList []string
	categoryListRaw := strings.Split(viper.GetString("ARXIV_CATEGORY_LIST"), ",")
	for _, category := range categoryListRaw {
		if category != "" {
			categoryList = append(categoryList, category)
		}
	}

	// arXiv's init URL list
	var initUrlList []string
	initUrlListRaw := strings.Split(viper.GetString("ARXIV_INIT_URL_LIST"), ",")
	for _, initUrl := range initUrlListRaw {
		if initUrl != "" {
			initUrlList = append(initUrlList, initUrl)
		}
	}

	// arXiv config
	Arxiv = &ArxivConfig{
		CategoryList:           categoryList,
		InitUrlList:            initUrlList,
		IsInsecureHttpAccepted: viper.GetBool("ARXIV_ACCEPT_INSECURE_HTTP"),
		RequestTimeout:         viper.GetDuration("ARXIV_REQUEST_TIMEOUT"),
		DuplicatedThreshold:    viper.GetInt("ARXIV_DUPLICATED_THRESHOLD"),
		MaxResults:             viper.GetInt("ARXIV_MAX_RESULTS"),
		SearchStart:            viper.GetInt("ARXIV_SEARCH_START"),
		SortBy:                 viper.GetString("ARXIV_SORT_BY"),
		SortOrder:              viper.GetString("ARXIV_SORT_ORDER"),
	}
}
