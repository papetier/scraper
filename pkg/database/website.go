package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"strings"
)

const websitesTable = "websites"

type Website struct {
	Id         ID       `db:"id"`
	DomainList []string `db:"-"`
	Name       string   `db:"name"`

	DomainListRaw string `db:"domain_list"`
}

func GetWebsites() ([]*Website, error) {
	var websiteList []*Website
	query := fmt.Sprintf(`SELECT id, name, domain_list FROM %s`, websitesTable)
	err := pgxscan.Select(context.Background(), dbConnection.Pool, &websiteList, query)
	if err != nil {
		return nil, err
	}

	for _, website := range websiteList {
		// allowed domains
		domainList := strings.Split(website.DomainListRaw, ",")
		for _, domain := range domainList {
			if domain != "" {
				website.DomainList = append(website.DomainList, domain)
			}
		}
	}

	return websiteList, nil
}
