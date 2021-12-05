package scraper

import (
	"github.com/gocolly/colly/v2"
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

	wc := collector.GetWebsiteCollector(website, colly.AllowURLRevisit())

	// websites specific collector settings
	switch website.Name {
	case "arXiv":
		err := arxiv.UpdateAndLoadCategories(website)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("categories successfully loaded --------- now starting scraper!")
		arxiv.SetupCollector(wc.Collector)
	}

	// on initial urls
	for _, initUrl := range website.InitUrlList {
		wc.AddUrl(initUrl)
	}

	// on category list for arXiv
	for _, category := range website.CategoryList {
		if website.Name == "arXiv" {
			arxiv.SearchCategory(wc, category)
		}
	}
}
