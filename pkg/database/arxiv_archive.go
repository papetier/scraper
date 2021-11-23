package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
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
			archiveName = archive.ArxivGroup.OriginalArxivGroupName
		}
		archive.OriginalArxivArchiveName = archiveName

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
		archive.OriginalArxivArchiveCode = archiveCode

		archiveValues = append(archiveValues, archiveCode, archiveName, archive.ArxivGroupId)
	}

	archivePlaceholder := generateInsertPlaceholder(len(arxivArchiveColumns[1:]), len(archiveList), 1)
	archivesQuery := "INSERT INTO " + arxivArchivesTable + " (" + strings.Join(arxivArchiveColumns[1:], ", ") + ") VALUES " + archivePlaceholder + " ON CONFLICT DO NOTHING RETURNING id"

	archiveRows, err := dbConnection.Pool.Query(context.Background(), archivesQuery, archiveValues...)
	defer archiveRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the arXiv's archives into the database: %w", err)
	}

	insertedArchiveCount := 0
	var insertedArchiveIdList []ID
	for archiveRows.Next() {
		err = archiveRows.Scan(&insertedArchiveIdList[insertedArchiveCount])
		if err != nil {
			return fmt.Errorf("scanning the arXiv's archive ids: %w", err)
		}
		insertedArchiveCount++
	}

	if insertedArchiveCount == len(archiveList) {
		for i, id := range insertedArchiveIdList {
			archiveList[i].Id = id
			updateArchiveReferenceInArxivCategories(archiveList[i])
		}
	} else {
		err = fetchArxivArchiveIds(archiveList)
		if err != nil {
			return fmt.Errorf("fetching the arXiv's archive ids: %w", err)
		}
	}

	return nil
}

func fetchArxivArchiveIds(archiveList []*ArxivArchive) error {
	query := "SELECT id, original_arxiv_archive_code FROM " + arxivArchivesTable
	var fetchedArchiveList []*ArxivArchive
	err := pgxscan.Select(context.Background(), dbConnection.Pool, &fetchedArchiveList, query)
	if err != nil {
		return fmt.Errorf("scanning the arxiv archive list: %w", err)
	}

	arxivArchiveIdMapByCode := make(map[string]ID)
	for _, arxivArchive := range fetchedArchiveList {
		arxivArchiveIdMapByCode[arxivArchive.OriginalArxivArchiveCode] = arxivArchive.Id
	}

	for _, archive := range archiveList {
		if archive.Id == 0 {
			archive.Id = arxivArchiveIdMapByCode[archive.OriginalArxivArchiveCode]
		}
		updateArchiveReferenceInArxivCategories(archive)
	}

	return nil
}

func updateArchiveReferenceInArxivCategories(archive *ArxivArchive) {
	for _, category := range archive.ArxivCategories {
		category.ArxivArchiveId = archive.Id
	}
}
