package scraper

import (
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const layoutParuvendu = "Le 02/01/2006 à 15:04"

var p *bluemonday.Policy
var location *time.Location

func setupHtmlSanitizer() {
	p = bluemonday.StripTagsPolicy()
}

func sanitize(html string) string {
	return p.Sanitize(html)
}

func setupTimeLocation() {
	var err error
	location, err = time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatal(err)
	}
}

func StandardizeStringWhitespaces(input string) string {
	return strings.Join(strings.Fields(input), " ")
}
