package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ArxivEprint struct {
	Id ID `db:"id"`

	ArxivId string `db:"arxiv_id"`

	Comment       *string                 `db:"comment"`
	Extra         *map[string]interface{} `db:"extra"`
	LatestVersion int                     `db:"latest_version"`
	PdfLink       *string                 `db:"pdf_link"`

	PublishedAt time.Time `db:"published_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	PaperId ID `db:"paper_id"`

	Paper                *Paper
	PrimaryArxivCategory *ArxivCategory
	OtherArxivCategories []*ArxivCategory
}

const (
	arxivEprintsTable                = "arxiv_eprints"
	arxivEprintsArxivCategoriesTable = "arxiv_eptrins_arxiv_categories"
)

var arxivEprintsColumns = []string{
	"id",
	"arxiv_id",
	"paper_id",
	"comment",
	"extra",
	"latest_version",
	"pdf_link",
	"published_at",
	"updated_at",
}

var arxivEprintsArxivCategoriesColumns = []string{
	"arxiv_eprint_id",
	"arxiv_category_id",
	"is_primary",
}

func (a *ArxivEprint) SaveWithPaperAuthorsAndCategories() error {
	log.Debugf("saving arXiv's eprint `%s` with related paper and authors", a.ArxivId)

	// prepare transaction
	tx, err := dbConnection.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// save authors w/ organisations
	err = saveAuthorsWithOrganisationsTx(tx, a.Paper.Authors)
	if err != nil {
		return fmt.Errorf("saving the authors with their organisations associated with the arXiv's eprint's `%s`: %w", a.ArxivId, err)
	}

	// save paper with author links (and author order)
	err = a.Paper.SaveWithAuthorsTx(tx)
	if err != nil {
		return fmt.Errorf("saving the paper associated with the arXiv's eprint `%s`: %w", a.ArxivId, err)
	}
	a.PaperId = a.Paper.Id

	// save arxiv_eprint with categories
	err = a.saveWithCategoriesTx(tx)
	if err != nil {
		return fmt.Errorf("saving the arXiv's eprint `%s` with categories: %w", a.ArxivId, err)
	}

	// commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("committing the transaction to save the arXiv's eprint `%s`, paper and authors: %w", a.ArxivId, err)
	}

	return nil
}

func (a *ArxivEprint) saveWithCategoriesTx(tx pgx.Tx) error {
	// save arxiv_eprint
	err := a.saveTx(tx)
	if err != nil {
		return fmt.Errorf("saving the arxiv_eprint: %w", err)
	}

	// save links arxiv_eprint/categories
	if a.PrimaryArxivCategory != nil || len(a.OtherArxivCategories) > 0 {
		err = a.saveCategoriesTx(tx)
		if err != nil {
			return fmt.Errorf("saving the arxiv_eprint_arxiv_categories: %w", err)
		}
	}

	return nil
}

func (a *ArxivEprint) saveTx(tx pgx.Tx) error {
	log.Debugf("saving the arXiv's eprint `%s`", a.ArxivId)

	arxivEprintPlaceholder := generateInsertPlaceholder(len(arxivEprintsColumns[1:]), 1, 1)
	arxivEprintsQuery := "INSERT INTO " + arxivEprintsTable + " (" + strings.Join(arxivEprintsColumns[1:], ", ") + ") VALUES " + arxivEprintPlaceholder + " RETURNING id"

	arxivEprintRow, err := tx.Query(context.Background(), arxivEprintsQuery, a.ArxivId, a.Paper.Id, a.Comment, a.Extra, a.LatestVersion, a.PdfLink, a.PublishedAt, a.UpdatedAt)
	defer arxivEprintRow.Close()
	if err != nil {
		return fmt.Errorf("inserting the arxiv_eprint into the database: %w", err)
	}

	for arxivEprintRow.Next() {
		err = arxivEprintRow.Scan(&a.Id)
		if err != nil {
			return fmt.Errorf("scanning the arxiv_eprint id: %w", err)
		}
	}

	return nil
}

func (a *ArxivEprint) saveCategoriesTx(tx pgx.Tx) error {
	log.Debug("saving the arxiv_eprints_arxiv_categories links")

	categoryCount := 0
	var categoryLinkValues []interface{}
	if a.PrimaryArxivCategory != nil {
		categoryLinkValues = append(categoryLinkValues, a.Id, a.PrimaryArxivCategory.Id, true)
		categoryCount++
	}
	for _, category := range a.OtherArxivCategories {
		categoryLinkValues = append(categoryLinkValues, a.Id, category.Id, false)
		categoryCount++
	}

	categoryLinkPlaceholder := generateInsertPlaceholder(len(arxivEprintsArxivCategoriesColumns), categoryCount, 1)
	authorLinksQuery := "INSERT INTO " + arxivEprintsArxivCategoriesTable + " (" + strings.Join(arxivEprintsArxivCategoriesColumns, ", ") + ") VALUES " + categoryLinkPlaceholder

	categoryLinkRows, err := tx.Query(context.Background(), authorLinksQuery, categoryLinkValues...)
	defer categoryLinkRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the arxiv_eprints_arxiv_categories links into the database: %w", err)
	}

	return nil
}
