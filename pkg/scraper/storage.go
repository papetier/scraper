package scraper

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/papetier/scraper/pkg/database"
	"log"
	"net/url"
	"strconv"
	"sync"
)

const visitedTable = "colly_storage_visited_pages"
const cookiesTable = "colly_storage_cookies"

// From https://github.com/zolamk/colly-postgres-storage/blob/master/colly/postgres/storage.go

// Storage implements a PostgreSQL storage backend for colly
type Storage struct {
	pool         *pgxpool.Pool
	VisitedTable string
	CookiesTable string
}

var once sync.Once

func prepareDB(s *Storage) {
	once.Do(func() {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (request_id text not null);", s.VisitedTable)
		_, err := s.pool.Exec(context.Background(), query)
		if err != nil {
			log.Fatal(err)
		}

		query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (host text not null, cookies text not null);", s.CookiesTable)
		_, err = s.pool.Exec(context.Background(), query)
		if err != nil {
			log.Fatal(err)
		}
	})
}

//Init initializes the PostgreSQL storage
func (s *Storage) Init() error {
	s.pool = database.GetPool()

	err := s.pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	prepareDB(s)

	return nil
}

// Visited implements colly/storage.Visited()
func (s *Storage) Visited(requestID uint64) error {
	var err error

	query := fmt.Sprintf(`INSERT INTO %s (request_id) VALUES($1);`, s.VisitedTable)

	_, err = s.pool.Exec(context.Background(), query, strconv.FormatUint(requestID, 10))

	return err
}

// IsVisited implements colly/storage.IsVisited()
func (s *Storage) IsVisited(requestID uint64) (bool, error) {

	var isVisited bool

	query := fmt.Sprintf(`SELECT EXISTS(SELECT request_id FROM %s WHERE request_id = $1)`, s.VisitedTable)

	err := s.pool.QueryRow(context.Background(), query, strconv.FormatUint(requestID, 10)).Scan(&isVisited)

	return isVisited, err
}

// Cookies implements colly/storage.Cookies()
func (s *Storage) Cookies(u *url.URL) string {

	var cookies string

	query := fmt.Sprintf(`SELECT cookies FROM %s WHERE host = $1;`, s.CookiesTable)

	s.pool.QueryRow(context.Background(), query, u.Host).Scan(&cookies)

	return cookies
}

// SetCookies implements colly/storage.SetCookies()
func (s *Storage) SetCookies(u *url.URL, cookies string) {

	query := fmt.Sprintf(`INSERT INTO %s (host, cookies) VALUES($1, $2);`, s.CookiesTable)

	s.pool.Exec(context.Background(), query, u.Host, cookies)
}
