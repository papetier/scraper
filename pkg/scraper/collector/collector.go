package collector

import (
	"github.com/gocolly/colly/v2"
	"github.com/papetier/scraper/pkg/database"
	"github.com/papetier/scraper/pkg/scraper/storage"
	log "github.com/sirupsen/logrus"
)

type WebsiteCollector struct {
	Website   *database.Website
	Collector *colly.Collector
}

func (wc *WebsiteCollector) AddUrl(url string) {
	err := wc.Collector.Visit(url)
	if err != nil {
		log.Errorf("%s collector encountered an error when fetching %s\nError: %v", wc.Website.Domain, url, err)
	}
}

func GetWebsiteCollector(website *database.Website, options ...colly.CollectorOption) *WebsiteCollector {
	// new colly collector
	collectorOptions := options
	collectorOptions = append(collectorOptions, colly.AllowedDomains(website.Domain))
	c := colly.NewCollector(collectorOptions...)

	// storage set up
	err := c.SetStorage(storage.DbStorage)
	if err != nil {
		log.Fatal(err)
	}

	// basic callbacks
	c.OnRequest(func(r *colly.Request) {
		log.Infof("fetching: %s", r.URL)
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Warnf("request URL: %v failed with response: %v\nError: %v", r.Request.URL.String(), r, err)
	})
	c.OnScraped(onScraped())

	return &WebsiteCollector{
		Website:   website,
		Collector: c,
	}
}

func onScraped() func(r *colly.Response) {
	return func(r *colly.Response) {
		log.Infof("finished scraping: %s", r.Request.URL.String())
	}
}
