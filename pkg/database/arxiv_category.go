package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
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

	updatedCategoryCount := 0
	for categoryRows.Next() {
		err = categoryRows.Scan(&categoryList[updatedCategoryCount].Id)
		if err != nil {
			return fmt.Errorf("scanning the arXiv's category ids: %w", err)
		}
		updatedCategoryCount++
	}

	if updatedCategoryCount < len(categoryList) {
		arxivArchiveIdMapByCode, err := getArxivCategoryIdMapByCode()
		if err != nil {
			return fmt.Errorf("fetching the group ids: %w", err)
		}

		for _, category := range categoryList {
			if category.Id == 0 {
				category.Id = arxivArchiveIdMapByCode[category.OriginalArxivCategoryCode]
			}
		}
	}

	return nil
}

func getArxivCategoryIdMapByCode() (map[string]ID, error) {
	query := "SELECT id, original_arxiv_category_code FROM " + arxivCategoriesTable
	var arxivCategoryList []*ArxivCategory
	err := pgxscan.Select(context.Background(), dbConnection.Pool, &arxivCategoryList, query)
	if err != nil {
		return nil, fmt.Errorf("scanning the arxiv category list: %w", err)
	}

	arxivCategoryIdMapByCode := make(map[string]ID)
	for _, arxivCategory := range arxivCategoryList {
		arxivCategoryIdMapByCode[arxivCategory.OriginalArxivCategoryCode] = arxivCategory.Id
	}

	return arxivCategoryIdMapByCode,nil
}
