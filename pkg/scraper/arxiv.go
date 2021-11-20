package scraper

import (
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func SetupArxivCollector(wc *WebsiteCollector) {
	wc.Collector.OnHTML("div#bloc_liste div.ergov3-annonce a[href^=\"/\"]", listParser(wc))
	wc.Collector.OnHTML("div#maindetail, div#sidebar-autodetail", adParser(wc))

}

func listParser(wc *WebsiteCollector) func(*colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		// TODO
		log.Debug("received element")
	}
}

func adParser(wc *WebsiteCollector) func(*colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		log.Debug("received element")
		// TODO
	}
}
