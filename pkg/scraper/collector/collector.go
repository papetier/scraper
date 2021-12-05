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
		log.Errorf("error visiting %s: %s", url, err)
	}
}

func GetWebsiteCollector(website *database.Website, options ...colly.CollectorOption) *WebsiteCollector {
	// new colly collector
	collectorOptions := options
	collectorOptions = append(collectorOptions, colly.AllowedDomains(website.DomainList...))
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
		log.WithField("collector", website.Name).Errorf("request URL: %v failed with HTTP status %v: %s", r.Request.URL.String(), r.StatusCode, err)
	})
	c.OnScraped(onScraped())

	return &WebsiteCollector{
		Website:   website,
		Collector: c,
	}
}

func onScraped() func(r *colly.Response) {
	return func(r *colly.Response) {
		log.Debugf("finished: %s", r.Request.URL.String())
	}
}
