package arxiv

import (
	"github.com/gocolly/colly/v2"
	"github.com/papetier/scraper/pkg/database"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	arxivAbstractUrl = arxivBaseUrl + "abs/"
	arxivPdfUrl      = arxivBaseUrl + "pdf/"
	arxivBaseUrl     = "http://arxiv.org/"
	arxivErrorTitle  = "Error"
)

func SetupCollector(c *colly.Collector) {
	c.OnXML("/feed/entry", entryParser)

}

func entryParser(e *colly.XMLElement) {
	title := e.ChildText("title")
	if title == arxivErrorTitle {
		handleErrorEntry(e)
		return
	}

	// initialise paper + arxiv eprint
	paper := &database.Paper{
		Title: title,
	}
	arxivEprint := &database.ArxivEprint{
		EPrint: &database.Eprint{
			Paper: paper,
		},
	}

	// parse id
	id := e.ChildText("id")
	idParsingResult := strings.Split(id, arxivAbstractUrl)
	if len(idParsingResult) < 2 {
		log.Errorf("unexpected arxiv id format: %s", id)
	} else {
		arxivId := idParsingResult[1]
		log.Debugf("parsing entry element %s", arxivId)
		arxivEprint.ArxivId = arxivId
	}

	// parse doi
	doi := e.ChildText("arxiv:doi")
	paper.Doi = doi

	// parse abstract
	abstract := e.ChildText("summary")
	paper.Abstract = abstract

	// parse journal_ref
	journalRef := e.ChildText("arxiv:journal_ref")
	paper.JournalRef = journalRef

	// TODO: parse year

	// parse comment
	comment := e.ChildText("arxiv:comment")
	if comment != "" {
		arxivEprint.Comment = &comment
	}

	// parse pdf_link (if different from default)
	pdfLink := e.ChildAttr("link[@title='pdf']", "href")
	if pdfLink != arxivPdfUrl+arxivEprint.ArxivId {
		arxivEprint.PdfLink = &pdfLink
	}

	// parse published
	publishedAtRaw := e.ChildText("published")
	publishedAt, err := time.Parse(time.RFC3339, publishedAtRaw)
	if err != nil {
		log.Errorf("parsing published date for %s : %s", arxivEprint.ArxivId, err)
	}
	arxivEprint.PublishedAt = publishedAt

	// parse updated
	updatedAtRaw := e.ChildText("updated")
	updatedAt, err := time.Parse(time.RFC3339, updatedAtRaw)
	if err != nil {
		log.Errorf("parsing updated date for %s : %s", arxivEprint.ArxivId, err)
	}
	arxivEprint.UpdatedAt = updatedAt

	// parse latest_version
	versionSplitResult := strings.Split(arxivEprint.ArxivId, "v")
	if len(versionSplitResult) < 2 {
		log.Errorf("unexpected arxiv id format for version parsing: %s", id)
	} else {
		latestVersionString := versionSplitResult[len(versionSplitResult)-1]
		latestVersion, err := strconv.Atoi(latestVersionString)
		if err != nil {
			log.Errorf("formatting latest version (%s): %s", latestVersionString, err)
		}
		arxivEprint.LatestVersion = latestVersion
	}

	// parse categories (with primary) + extra categories
	var otherArxivCategories []*database.ArxivCategory
	var extraCategories []string
	primaryCategoryCode := e.ChildAttr("arxiv:primary_category", "term")
	categoryCodeList := e.ChildAttrs("category", "term")
	for _, categoryCode := range categoryCodeList {
		if arxivCategory, exists := categoriesByCodeMap[categoryCode]; exists {
			if categoryCode == primaryCategoryCode {
				arxivEprint.PrimaryArxivCategory = arxivCategory
			} else {
				otherArxivCategories = append(otherArxivCategories, arxivCategory)
			}
		} else {
			extraCategories = append(extraCategories, categoryCode)
		}
	}
	arxivEprint.OtherArxivCategories = otherArxivCategories

	// parse extra (atm: only extra categories)
	if len(extraCategories) > 0 {
		arxivEprint.Extra = &map[string]interface{}{
			"categories": extraCategories,
		}
	}

	// TODO: save arxivEprint and co
}

func handleErrorEntry(e *colly.XMLElement) {
	// TODO
	log.Infof("an error occurred")
}
