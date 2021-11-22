package scraper

import (
	"github.com/papetier/scraper/pkg/database"
	"github.com/papetier/scraper/pkg/scraper/arxiv"
	"github.com/papetier/scraper/pkg/scraper/collector"
	"github.com/papetier/scraper/pkg/scraper/storage"
	log "github.com/sirupsen/logrus"
	"sync"
)


func Setup() {
	storage.SetupSDBStorage()
}

func ScrapeWebsites(websiteList []*database.Website) {
	var wg sync.WaitGroup

	for _, website := range websiteList {
		wg.Add(1)
		go scrape(website, &wg)
	}
	wg.Wait()
}

func scrape(website *database.Website, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Infof("Scraping %s...", website.Name)

	wc := collector.GetWebsiteCollector(website)

	// websites specific collector settings
	switch website.Name {
	case "arXiv":
		err := arxiv.UpdateAndLoadCategories(website)
		if err != nil {
			log.Fatal(err)
		}
		arxiv.SetupCollector(wc.Collector)
	}

	// run the collector
	for _, initUrl := range website.InitUrlList {
		wc.AddUrl(initUrl)
	}
}
