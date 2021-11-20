package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/papetier/scraper/pkg/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

type WebsiteCollector struct {
	Website   *database.Website
	Collector *colly.Collector
}

const envPrefix = "SCRAPE_"

func ScrapeAllEnabled(websiteList []*database.Website) {
	var wg sync.WaitGroup

	storage := &Storage{
		VisitedTable: visitedTable,
		CookiesTable: cookiesTable,
	}

	// TODO: from config
	for _, website := range websiteList {

		// check if this website should be scraped
		envKey := envPrefix + strings.ReplaceAll(strings.ToUpper(website.Name), " ", "_")
		shouldScrapeWebsite := viper.GetBool(envKey)

		// add it to the wait group
		if shouldScrapeWebsite {
			wg.Add(1)
			go scrapeApi(website, storage, &wg)
		}
	}
	wg.Wait()
}

func scrapeApi(website *database.Website, storage *Storage, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Scraping %s...", website.Name)

	// get base config
	c := getCollector(website)

	err := c.SetStorage(storage)
	if err != nil {
		log.Fatal(err)
	}

	// instantiate corresponding WebsiteCollector
	wc := &WebsiteCollector{
		Website:   website,
		Collector: c,
	}

	// websites specific collector settings
	switch website.Name {
	case "arXiv":
		SetupArxivCollector(wc)
	}

	// run the collector
	appendUrlToCollector(wc, website.InitUrl)
}

func getCollector(website *database.Website) *colly.Collector {
	// TODO: async/collector.Wait
	c := colly.NewCollector(
		colly.AllowedDomains(website.Domain),
		//colly.AllowURLRevisit(),
	)

	// basic callbacks
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %v failed with response: %v\nError: %v", r.Request.URL.String(), r, err)
	})

	c.OnScraped(onScraped())

	return c
}

func appendUrlToCollector(wc *WebsiteCollector, url string) {
	err := wc.Collector.Visit(url)
	if err != nil {
		log.Errorf("%s collector encountered an error when visiting %s\nError: %v", wc.Website.Domain, url, err)
	}
}

func onScraped() func(r *colly.Response) {
	return func(r *colly.Response) {
		log.Printf("> Finished scraping %s", r.Request.URL.String())
	}
}
