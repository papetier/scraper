package arxiv

import (
	"fmt"
	"github.com/papetier/scraper/pkg/scraper/collector"
	log "github.com/sirupsen/logrus"
)

const (
	arxivBaseSearchUrl = "http://export.arxiv.org/api/query?search_query="
	arxivQueryPattern  = "%scat:%s&start=%d&max_results=%d&sortBy=%s&sortOrder=%s"
	defaultMaxResults  = 1000
	defaultSortBy      = "submittedDate"
	defaultSortOrder   = "ascending"
)

func SearchCategory(wc *collector.WebsiteCollector, categoryCode string) {
	category, present := categoriesByCodeMap[categoryCode]
	if !present {
		log.Errorf("unknown arXiv category code: %s", categoryCode)
	}

	start := 0
	for duplicatedPaperCounter < 3 {
		queryString := fmt.Sprintf(arxivQueryPattern, arxivBaseSearchUrl, category.OriginalArxivCategoryCode, start, defaultMaxResults, defaultSortBy, defaultSortOrder)
		wc.AddUrl(queryString)
		start += defaultMaxResults
	}

	log.Infof("stopping search scraper with %d duplicated papers", duplicatedPaperCounter)
}
