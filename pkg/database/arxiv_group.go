package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ArxivGroup struct {
	Id                     ID     `db:"id"`
	OriginalArxivGroupName string `db:"original_arxiv_group_name"`

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
		groupValues = append(groupValues, group.OriginalArxivGroupName)
	}

	groupPlaceholder := generateInsertPlaceholder(len(arxivGroupColumns[1:]), len(groupList), 1)
	groupsQuery := "INSERT INTO " + arxivGroupsTable + " (" + strings.Join(arxivGroupColumns[1:], ", ") + ") VALUES " + groupPlaceholder + " ON CONFLICT DO NOTHING RETURNING id"

	groupRows, err := dbConnection.Pool.Query(context.Background(), groupsQuery, groupValues...)
	defer groupRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the arXiv's groups into the database: %w", err)
	}

	updatedGroupCount := 0
	for groupRows.Next() {
		err = groupRows.Scan(&groupList[updatedGroupCount].Id)
		if err != nil {
			return fmt.Errorf("scanning the arXiv's group ids: %w", err)
		}

		// update archives with the group ids
		updateGroupReferenceInArxivArchive(groupList[updatedGroupCount])

		updatedGroupCount++
	}

	if updatedGroupCount < len(groupList) {
		arxivGroupIdMapByName, err := getArxivGroupIdMapByName()
		if err != nil {
			return fmt.Errorf("fetching the group ids: %w", err)
		}

		for i, group := range groupList {
			if group.Id == 0 {
				group.Id = arxivGroupIdMapByName[group.OriginalArxivGroupName]
			}
			// update archives with the group ids
			updateGroupReferenceInArxivArchive(groupList[i])
		}
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

func getArxivGroupIdMapByName() (map[string]ID, error) {
	query := "SELECT id, original_arxiv_group_name FROM " + arxivGroupsTable
	var arxivGroupList []*ArxivGroup
	err := pgxscan.Select(context.Background(), dbConnection.Pool, &arxivGroupList, query)
	if err != nil {
		return nil, fmt.Errorf("scanning the arxiv group list: %w", err)
	}

	arxivGroupIdMapByName := make(map[string]ID)
	for _, arxivGroup := range arxivGroupList {
		arxivGroupIdMapByName[arxivGroup.OriginalArxivGroupName] = arxivGroup.Id
	}

	return arxivGroupIdMapByName,nil
}

func updateGroupReferenceInArxivArchive(group *ArxivGroup) {
	for _, archive := range group.ArxivArchives {
		archive.ArxivGroupId = group.Id
	}
}
