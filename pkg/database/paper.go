package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Paper struct {
	Id ID `db:"id"`

	Doi        *string `db:"doi"`
	JournalRef *string `db:"journal_ref"`

	Abstract string `db:"abstract"`
	Title    string `db:"title"`
	Year     *int   `db:"year"`

	Authors []*Author
}

const (
	papersTable        = "papers"
	papersAuthorsTable = "papers_authors"
)

var papersColumns = []string{
	"id",
	"doi",
	"journal_ref",
	"abstract",
	"title",
	"year",
}

var papersAuthorsColumns = []string{
	"paper_id",
	"author_id",
	"author_order",
}

func (p *Paper) SaveWithAuthorsTx(tx pgx.Tx) error {
	log.Debug("saving the paper with its authors")

	err := p.saveTx(tx)
	if err != nil {
		return fmt.Errorf("saving the paper: %w", err)
	}

	err = p.saveAuthorsTx(tx)
	if err != nil {
		return fmt.Errorf("saving the papers_authors links: %w", err)
	}

	return nil
}

func (p *Paper) saveTx(tx pgx.Tx) error {
	log.Debugf("saving paper %s", p.Title)

	paperPlaceholder := generateInsertPlaceholder(len(papersColumns[1:]), 1, 1)
	papersQuery := "INSERT INTO " + papersTable + " (" + strings.Join(papersColumns[1:], ", ") + ") VALUES " + paperPlaceholder + " RETURNING id"

	paperRow, err := tx.Query(context.Background(), papersQuery, p.Doi, p.JournalRef, p.Abstract, p.Title, p.Year)
	defer paperRow.Close()
	if err != nil {
		return fmt.Errorf("inserting the paper into the database: %w", err)
	}

	for paperRow.Next() {
		err = paperRow.Scan(&p.Id)
		if err != nil {
			return fmt.Errorf("scanning the paper id: %w", err)
		}
	}

	return nil
}

func (p *Paper) saveAuthorsTx(tx pgx.Tx) error {
	log.Debug("saving the papers_authors links")

	var authorLinkValues []interface{}
	for order, author := range p.Authors {
		authorLinkValues = append(authorLinkValues, p.Id, author.Id, order)
	}

	authorLinkPlaceholder := generateInsertPlaceholder(len(papersAuthorsColumns), len(p.Authors), 1)
	authorLinksQuery := "INSERT INTO " + papersAuthorsTable + " (" + strings.Join(papersAuthorsColumns, ", ") + ") VALUES " + authorLinkPlaceholder

	authorLinkRows, err := tx.Query(context.Background(), authorLinksQuery, authorLinkValues...)
	defer authorLinkRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the paper_authors links into the database: %w", err)
	}

	return nil
}
