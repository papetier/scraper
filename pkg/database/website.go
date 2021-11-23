package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"strings"
)

const websitesTable = "websites"

type Website struct {
	Id          ID       `db:"id"`
	DomainList  []string `db:"-"`
	Name        string   `db:"name"`
	InitUrlList []string `db:"-"`

	DomainListRaw  string `db:"domain_list"`
	InitUrlListRaw string `db:"init_url_list"`
}

func GetWebsites() ([]*Website, error) {
	var websiteList []*Website
	query := fmt.Sprintf(`SELECT id, name, domain_list, init_url_list FROM %s`, websitesTable)
	err := pgxscan.Select(context.Background(), dbConnection.Pool, &websiteList, query)
	if err != nil {
		return nil, err
	}

	for _, website := range websiteList {
		website.InitUrlList = strings.Split(website.InitUrlListRaw, ",")
		website.DomainList = strings.Split(website.DomainListRaw, ",")
	}

	return websiteList, nil
}
