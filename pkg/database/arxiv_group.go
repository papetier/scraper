package database

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ArxivGroup struct {
	Id                ID     `db:"id"`
	OriginalGroupName string `db:"original_group_name"`

	ArxivArchives []*ArxivArchive
}

const arxivGroupsTable = "arxiv_groups"

var arxivGroupColumns = []string{
	"id",
	"original_arxiv_group_name",
}

func SaveArxivGroupsArchivesAndCategories(groupList []*ArxivGroup) error {
	log.Debug("saving arXiv's groups, archives and categories")

	var groupValues []interface{}
	for _, group := range groupList {
		groupValues = append(groupValues, group.OriginalGroupName)
	}

	groupPlaceholder := generateInsertPlaceholder(len(arxivGroupColumns[1:]), len(groupList), 1)
	groupsQuery := "INSERT INTO " + arxivGroupsTable + " (" + strings.Join(arxivGroupColumns[1:], ", ") + ") VALUES " + groupPlaceholder + " ON CONFLICT DO NOTHING RETURNING id"

	groupRows, err := dbConnection.Pool.Query(context.Background(), groupsQuery, groupValues...)
	defer groupRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the arXiv's groups into the database: %w", err)
	}

	i := 0
	for groupRows.Next() {
		err = groupRows.Scan(&groupList[i].Id)
		if err != nil {
			return fmt.Errorf("scanning the arXiv's group ids: %w", err)
		}

		// update archives with the group ids
		for _, archive := range groupList[i].ArxivArchives {
			archive.ArxivGroupId = groupList[i].Id
		}

		i++
	}

	// prepare archiveList + categoryList
	var archiveList []*ArxivArchive
	var categoryList []*ArxivCategory
	for _, group := range groupList {
		archiveList = append(archiveList, group.ArxivArchives...)
		for _, archive := range group.ArxivArchives {
			categoryList = append(categoryList, archive.ArxivCategories...)
		}
	}

	// save the archives
	err = saveArxivArchives(archiveList)
	if err != nil {
		return fmt.Errorf("saving the arXiv's archives: %w", err)
	}

	// save the categories
	err = saveArxivCategories(categoryList)
	if err != nil {
		return fmt.Errorf("saving the arXiv's categories: %w", err)
	}

	return nil
}
