package arxiv

import (
	"github.com/gocolly/colly/v2"
)

func SetupCollector(c *colly.Collector) {
	// TODO
	//c.OnXML("div#bloc_liste div.ergov3-annonce a[href^=\"/\"]", listParser(c))
	//c.OnXML("div#maindetail, div#sidebar-autodetail", entryParser(c))

}

//func listParser(c *colly.Collector) func(element *colly.XMLElement) {
//	return func(e *colly.XMLElement) {
//		// TODO
//		log.Debug("received element")
//	}
//}
//
//func entryParser(c *colly.Collector) func(element *colly.XMLElement) {
//	return func(e *colly.XMLElement) {
//		log.Debug("received element")
//		// TODO
//	}
//}

