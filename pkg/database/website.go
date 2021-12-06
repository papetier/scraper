package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"strings"
)

const websitesTable = "websites"

type Website struct {
	Id           ID       `db:"id"`
	CategoryList []string `db:"-"`
	DomainList   []string `db:"-"`
	Name         string   `db:"name"`
	InitUrlList  []string `db:"-"`

	CategoryListRaw string `db:"category_list"`
	DomainListRaw   string `db:"domain_list"`
	InitUrlListRaw  string `db:"init_url_list"`
}

func GetWebsites() ([]*Website, error) {
	var websiteList []*Website
	query := fmt.Sprintf(`SELECT id, name, category_list, domain_list, init_url_list FROM %s`, websitesTable)
	err := pgxscan.Select(context.Background(), dbConnection.Pool, &websiteList, query)
	if err != nil {
		return nil, err
	}

	for _, website := range websiteList {
		// categories to scrape
		categoryList := strings.Split(website.CategoryListRaw, ",")
		for _, category := range categoryList {
			if category != "" {
				website.CategoryList = append(website.CategoryList, category)
			}
		}

		// allowed domains
		domainList := strings.Split(website.DomainListRaw, ",")
		for _, domain := range domainList {
			if domain != "" {
				website.DomainList = append(website.DomainList, domain)
			}
		}

		// initial urls
		initUrlList := strings.Split(website.InitUrlListRaw, ",")
		for _, initUrl := range initUrlList {
			if initUrl != "" {
				website.InitUrlList = append(website.InitUrlList, initUrl)
			}
		}
	}

	return websiteList, nil
}
