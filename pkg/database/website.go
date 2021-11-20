package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
)

const websitesTable = "websites"

type Website struct {
	Id      ID
	Name    string
	Domain  string
	InitUrl string
}

func GetWebsites() ([]*Website, error) {
	var websites []*Website
	query := fmt.Sprintf(`SELECT id, name, domain, init_url FROM %s`, websitesTable)
	err := pgxscan.Select(context.Background(), dbConnection.Pool, &websites, query)
	if err != nil {
		return nil, err
	}

	return websites, nil
}
