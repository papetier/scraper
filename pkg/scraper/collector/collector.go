package collector

import (
	"crypto/tls"
	"github.com/gocolly/colly/v2"
	"github.com/papetier/scraper/pkg/config"
	"github.com/papetier/scraper/pkg/database"
	"github.com/papetier/scraper/pkg/scraper/storage"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
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

	// http settings
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Scraper.IsInsecureHttpAccepted},
		DialContext: (&net.Dialer{
			Timeout: config.Scraper.RequestTimeout,
		}).DialContext,
	})
	c.SetRequestTimeout(config.Scraper.RequestTimeout)

	// slow down colly to avoid saturating arXiv
	// following https://arxiv.org/help/api/tou#limitations
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*arxiv*",
		Parallelism: 1,
		Delay:       3 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	// storage set up
	err = c.SetStorage(storage.DbStorage)
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
