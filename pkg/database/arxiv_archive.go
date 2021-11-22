package database

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ArxivArchive struct {
	Id                       ID     `db:"id"`
	OriginalArxivArchiveCode string `db:"original_arxiv_archive_code"`
	OriginalArxivArchiveName string `db:"original_arxiv_archive_name"`

	ArxivGroupId ID `db:"arxiv_group_id"`

	ArxivGroup      *ArxivGroup
	ArxivCategories []*ArxivCategory
}

const arxivArchivesTable = "arxiv_archives"

var arxivArchiveColumns = []string{
	"id",
	"original_arxiv_archive_code",
	"original_arxiv_archive_name",
	"arxiv_group_id",
}

func saveArxivArchives(archiveList []*ArxivArchive) error {
	log.Debug("saving arXiv's archives")

	var archiveValues []interface{}
	for _, archive := range archiveList {
		archiveName := archive.OriginalArxivArchiveName
		if archiveName == "" {
			// defaults to group's name
			archiveName = archive.ArxivGroup.OriginalGroupName
		}

		archiveCode := archive.OriginalArxivArchiveCode
		if archiveCode == "" {
			// defaults to group's code (deduced from category's first part code)
			if len(archive.ArxivCategories) < 1 {
				log.Fatal("an unnamed arXiv's archive requires at least 1 category")
			}
			result := strings.Split(archive.ArxivCategories[0].OriginalArxivCategoryCode, ".")
			if len(result) < 1 || result[0] == "" {
				log.Fatalf("invalid category name: %s", archive.ArxivCategories[0].OriginalArxivCategoryCode)
			}
			archiveCode = result[0]
		}

		archiveValues = append(archiveValues, archiveCode, archiveName, archive.ArxivGroupId)
	}

	archivePlaceholder := generateInsertPlaceholder(len(arxivArchiveColumns[1:]), len(archiveList), 1)
	archivesQuery := "INSERT INTO " + arxivArchivesTable + " (" + strings.Join(arxivArchiveColumns[1:], ", ") + ") VALUES " + archivePlaceholder + " ON CONFLICT DO NOTHING RETURNING id"

	archiveRows, err := dbConnection.Pool.Query(context.Background(), archivesQuery, archiveValues...)
	defer archiveRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the arXiv's archives into the database: %w", err)
	}

	i := 0
	for archiveRows.Next() {
		err = archiveRows.Scan(&archiveList[i].Id)
		if err != nil {
			return fmt.Errorf("scanning the arXiv's archive ids: %w", err)
		}

		// update categories with the archive ids
		for _, category := range archiveList[i].ArxivCategories {
			category.ArxivArchiveId = archiveList[i].Id
		}

		i++
	}

	return nil
}
