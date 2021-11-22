package database

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ArxivCategory struct {
	Id                               ID     `db:"id"`
	OriginalArxivCategoryCode        string `db:"original_arxiv_category_code"`
	OriginalArxivCategoryDescription string `db:"original_arxiv_category_description"`
	OriginalArxivCategoryName        string `db:"original_arxiv_category_name"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	ArxivArchiveId ID `db:"arxiv_archive_id"`

	IsPrimary *bool

	ArxivArchive *ArxivArchive
}

const arxivCategoriesTable = "arxiv_categories"

var arxivCategoryColumns = []string{
	"id",
	"original_arxiv_category_code",
	"original_arxiv_category_description",
	"original_arxiv_category_name",
	"arxiv_archive_id",
}

func saveArxivCategories(categoryList []*ArxivCategory) error {
	log.Debug("saving arXiv's categories")

	var categoryValues []interface{}
	for _, category := range categoryList {
		categoryValues = append(categoryValues, category.OriginalArxivCategoryCode, category.OriginalArxivCategoryDescription, category.OriginalArxivCategoryName, category.ArxivArchiveId)
	}

	categoryPlaceholder := generateInsertPlaceholder(len(arxivCategoryColumns[1:]), len(categoryList), 1)
	categoriesQuery := "INSERT INTO " + arxivCategoriesTable + " (" + strings.Join(arxivCategoryColumns[1:], ", ") + ") VALUES " + categoryPlaceholder + " ON CONFLICT DO NOTHING RETURNING id"

	categoryRows, err := dbConnection.Pool.Query(context.Background(), categoriesQuery, categoryValues...)
	defer categoryRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the arXiv's categories into the database: %w", err)
	}

	i := 0
	for categoryRows.Next() {
		err = categoryRows.Scan(&categoryList[i].Id)
		if err != nil {
			return fmt.Errorf("scanning the arXiv's category ids: %w", err)
		}
		i++
	}

	return nil
}
