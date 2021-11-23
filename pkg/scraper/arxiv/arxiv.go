package arxiv

import (
	"github.com/gocolly/colly/v2"
	"github.com/papetier/scraper/pkg/database"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	arxivAbstractUrl = arxivBaseUrl + "abs/"
	arxivBaseUrl     = "http://arxiv.org/"
	arxivErrorTitle  = "Error"
)

func SetupCollector(c *colly.Collector) {
	c.OnXML("/feed/entry", entryParser)

}

func entryParser(e *colly.XMLElement) {
	title := e.ChildText("//title")
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
	id := e.ChildText("//id")
	idParsingResult := strings.Split(id, arxivAbstractUrl)
	if len(idParsingResult) < 2 {
		log.Errorf("unexpected arxiv id format: %s", id)
	} else {
		arxivId := idParsingResult[1]
		log.Debugf("parsing entry element %s", arxivId)
		arxivEprint.ArxivId = arxivId
	}

	// TODO: parse doi
	// TODO: parse abstract
	// TODO: parse year
	// TODO: parse journal_ref

	// TODO: parse published
	// TODO: parse updated
	// TODO: parse categories (with primary)
	// TODO: parse extra
	// TODO: parse latest_version
	// TODO: parse pdf_link (if different from default)

	// TODO: save arxivEprint and co
	log.Infof("new arxivEprint: %#v", arxivEprint)
	log.Infof("new paper: %#v", paper)
}

func handleErrorEntry(e *colly.XMLElement) {
	// TODO
	log.Infof("an error occurred")
}
