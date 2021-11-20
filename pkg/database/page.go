package database

import (
	"context"
	"fmt"
"net/url"
"time"
)

const pagesTable = "pages"

type Page struct {
	Id        ID
	Url       *url.URL
	VisitedAt time.Time
	Website   *Website
}

func SavePage(page *Page) error {
	query := fmt.Sprintf(`
INSERT INTO %s (url, visited_at, website_id)
VALUES ($1, NOW(), $2)
ON CONFLICT ON CONSTRAINT unique_pages_url
DO 
   UPDATE SET visited_at = NOW()
RETURNING id
`, pagesTable)
	row := dbConnection.Pool.QueryRow(context.Background(), query, page.Url.String(), page.Website.Id)
	err := row.Scan(&page.Id)
	if err != nil {
		return err
	}
	return nil
}

