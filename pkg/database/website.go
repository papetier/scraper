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
		website.CategoryList = strings.Split(website.CategoryListRaw, ",")
		website.DomainList = strings.Split(website.DomainListRaw, ",")
		website.InitUrlList = strings.Split(website.InitUrlListRaw, ",")
	}

	return websiteList, nil
}
