package arxiv

import (
	"github.com/antchfx/xmlquery"
	"github.com/gocolly/colly/v2"
	"github.com/papetier/scraper/pkg/config"
	"github.com/papetier/scraper/pkg/database"
	"github.com/papetier/scraper/pkg/scraper/collector"
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
	c.OnXML("/feed", feedParser)
	c.OnXML("/feed/entry", entryParser)
}

func VisitInitUrlList(wc *collector.WebsiteCollector) {
	for _, initUrl := range config.Arxiv.InitUrlList {
		wc.AddUrl(initUrl)
	}
}

func feedParser(e *colly.XMLElement) {
	// get category code
	feedTitle := e.ChildText("title")
	categoryCode := getCategoryCodeFromSearchFeedTitle(feedTitle)
	if categoryCode == nil {
		log.Errorf("no category found in feed title %s", feedTitle)
		return
	}

	// get first entry if exists
	firstEntry := xmlquery.FindOne(e.DOM.(*xmlquery.Node), "entry")

	// set isLastResultEmpty for this category
	if firstEntry == nil {
		isLastResultEmptyByCategoryCode[*categoryCode] = true
	} else {
		isLastResultEmptyByCategoryCode[*categoryCode] = false
	}
}

func entryParser(e *colly.XMLElement) {
	title := strings.TrimSpace(e.ChildText("title"))
	if title == arxivErrorTitle {
		handleErrorEntry(e)
		return
	}

	// get category code
	canonicalCategoryCode := getCategoryCodeFromSearchFeedTitle(title)

	// initialise paper + arxiv eprint
	paper := &database.Paper{
		Title: title,
	}
	arxivEprint := &database.ArxivEprint{
		Paper: paper,
	}

	// parse id
	id := strings.TrimSpace(e.ChildText("id"))
	idParsingResult := strings.Split(id, arxivAbstractUrl)
	if len(idParsingResult) < 2 {
		log.Errorf("unexpected arxiv id format: %s", id)
		handleErrorEntry(e)
		return
	} else {
		arxivId := idParsingResult[1]
		log.Debugf("parsing entry element %s", arxivId)
		arxivEprint.ArxivId = arxivId
	}

	// parse doi
	doi := strings.TrimSpace(e.ChildText("arxiv:doi"))
	if doi != "" {
		paper.Doi = &doi
	}

	// parse abstract
	abstract := strings.TrimSpace(e.ChildText("summary"))
	paper.Abstract = abstract

	// parse journal_ref
	journalRef := strings.TrimSpace(e.ChildText("arxiv:journal_ref"))
	if journalRef != "" {
		paper.JournalRef = &journalRef
	}

	// TODO: parse year

	// parse authors
	var authorList []*database.Author
	authorIndex := 1
	for {
		authorName, authorAffiliation := getAuthorNameAndAffiliation(e, authorIndex)
		if authorName == "" {
			break
		}

		author := &database.Author{
			FullName: authorName,
		}

		if authorAffiliation != "" {
			organisation := &database.Organisation{
				Name: authorAffiliation,
			}
			author.Organisations = []*database.Organisation{organisation}
		}

		authorList = append(authorList, author)
		authorIndex++
	}
	if len(authorList) > 0 {
		paper.Authors = authorList
	}

	// parse comment
	comment := strings.TrimSpace(e.ChildText("arxiv:comment"))
	if comment != "" {
		arxivEprint.Comment = &comment
	}

	// parse pdf_link (if different from default)
	pdfLink := strings.TrimSpace(e.ChildAttr("link[@title='pdf']", "href"))
	if pdfLink != arxivPdfUrl+arxivEprint.ArxivId {
		arxivEprint.PdfLink = &pdfLink
	}

	// parse published
	publishedAtRaw := strings.TrimSpace(e.ChildText("published"))
	publishedAt, err := time.Parse(time.RFC3339, publishedAtRaw)
	if err != nil {
		log.Errorf("parsing published date for %s : %s", arxivEprint.ArxivId, err)
	}
	arxivEprint.PublishedAt = publishedAt

	// parse updated
	updatedAtRaw := strings.TrimSpace(e.ChildText("updated"))
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
	primaryCategoryCode := strings.TrimSpace(e.ChildAttr("arxiv:primary_category", "term"))
	categoryCodeList := e.ChildAttrs("category", "term")
	for _, categoryCodeRaw := range categoryCodeList {
		categoryCode := strings.TrimSpace(categoryCodeRaw)
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

	isDuplicate, err := arxivEprint.SaveWithPaperAuthorsAndCategories()
	if canonicalCategoryCode != nil {
		if isDuplicate {
			duplicatedPaperCounterByCategoryCode[*canonicalCategoryCode]++
			log.Warnf("arXiv's eprint %s was already saved, skipping", arxivEprint.ArxivId)
		} else {
			duplicatedPaperCounterByCategoryCode[*canonicalCategoryCode] = 0
		}
	}
	if err != nil {
		log.Errorf("saving the arXiv's eprint: %s", err)
	}
}

func getAuthorNameAndAffiliation(e *colly.XMLElement, authorIndex int) (string, string) {
	xpathQuery := "author[" + strconv.Itoa(authorIndex) + "]"
	name := e.ChildText(xpathQuery + "/name")
	affiliation := e.ChildText(xpathQuery + "/arxiv:affiliation")
	return strings.TrimSpace(name), strings.TrimSpace(affiliation)
}

func handleErrorEntry(e *colly.XMLElement) {
	summary := e.ChildText("summary")
	log.Errorf("the query URL was malformed and the arXiv's API answered with an error: %s", summary)
}
