package arxiv

import (
	"fmt"
	"github.com/papetier/scraper/pkg/scraper/collector"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

const (
	arxivBaseSearchUrl         = "http://export.arxiv.org/api/query?search_query="
	arxivQueryPattern          = "%scat:%s&start=%d&max_results=%d&sortBy=%s&sortOrder=%s"
	defaultMaxResults          = 1000
	defaultSortBy              = "submittedDate"
	defaultSortOrder           = "ascending"
	searchQueryCategoryPattern = `cat:(.+)`
	searchQueryTitlePattern    = `(.*): search_query=(.*)&id_list=(.*)&start=(\d+)&max_results=(\d+)`
)

var searchQueryTitleRegex = regexp.MustCompile(searchQueryTitlePattern)
var searchQueryCategoryRegex = regexp.MustCompile(searchQueryCategoryPattern)

var duplicatedPaperCounterByCategoryCode map[string]int
var isLastResultEmptyByCategoryCode map[string]bool

func SearchCategoryList(wc *collector.WebsiteCollector) {
	// prepare tracker maps
	duplicatedPaperCounterByCategoryCode = make(map[string]int)
	isLastResultEmptyByCategoryCode = make(map[string]bool)

	for _, category := range wc.Website.CategoryList {
		if wc.Website.Name == "arXiv" {
			searchCategory(wc, category)
		}
	}
}

func searchCategory(wc *collector.WebsiteCollector, categoryCode string) {
	category, present := categoriesByCodeMap[categoryCode]
	if !present {
		log.Errorf("unknown arXiv category code: %s", categoryCode)
		return
	}

	start := 0
	for duplicatedPaperCounterByCategoryCode[categoryCode] < 3 && !isLastResultEmptyByCategoryCode[categoryCode] {
		queryString := fmt.Sprintf(arxivQueryPattern, arxivBaseSearchUrl, category.OriginalArxivCategoryCode, start, defaultMaxResults, defaultSortBy, defaultSortOrder)
		wc.AddUrl(queryString)
		start += defaultMaxResults
	}

	if isLastResultEmptyByCategoryCode[categoryCode] {
		log.Infof("last visited URL had an empty feed - stopping scraper search for category %s", categoryCode)
	} else {
		log.Infof("last visited URL resulted in %d duplicated entries - stopping scraper search for category %s", duplicatedPaperCounterByCategoryCode[categoryCode], categoryCode)
	}
}

func getCategoryCodeFromSearchFeedTitle(title string) *string {
	// extract canonical query
	queryResult := searchQueryTitleRegex.FindStringSubmatch(title)

	if len(queryResult) < 3 {
		return nil
	}
	sq := queryResult[2]

	// extract categories from search query parameters
	categoryResult := searchQueryCategoryRegex.FindStringSubmatch(sq)
	if len(categoryResult) < 2 {
		return nil
	}
	rawCategory := categoryResult[1]

	// extract first category
	categoryCodeList := strings.Split(rawCategory, ",")
	if len(categoryCodeList) < 1 {
		return nil
	}
	return &categoryCodeList[0]
}
