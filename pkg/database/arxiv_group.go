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

var arxivGroupsColumns = []string{
	"id",
	"original_arxiv_group_name",
}

func SaveArxivGroupsArchivesAndCategories(groupList []*ArxivGroup) error {
	log.Debug("saving arXiv's groups, archives and categories")

	var groupValues []interface{}
	for _, group := range groupList {
		groupValues = append(groupValues, group.OriginalArxivGroupName)
	}

	groupPlaceholder := generateInsertPlaceholder(len(arxivGroupsColumns[1:]), len(groupList), 1)
	groupsQuery := "INSERT INTO " + arxivGroupsTable + " (" + strings.Join(arxivGroupsColumns[1:], ", ") + ") VALUES " + groupPlaceholder + " ON CONFLICT DO NOTHING RETURNING id"

	groupRows, err := dbConnection.Pool.Query(context.Background(), groupsQuery, groupValues...)
	defer groupRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the arXiv's groups into the database: %w", err)
	}

	insertedGroupCount := 0
	var insertedGroupIdList []ID
	for groupRows.Next() {
		var id ID
		err = groupRows.Scan(&id)
		insertedGroupIdList = append(insertedGroupIdList, id)
		if err != nil {
			return fmt.Errorf("scanning the arXiv's group ids: %w", err)
		}
		insertedGroupCount++
	}

	if insertedGroupCount == len(groupList) {
		for i, id := range insertedGroupIdList {
			groupList[i].Id = id
			updateGroupReferenceInArxivArchive(groupList[i])
		}
	} else {
		err = fetchAndUpdateArxivGroupIds(groupList)
		if err != nil {
			return fmt.Errorf("fetching the arXiv's group ids: %w", err)
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

func fetchAndUpdateArxivGroupIds(groupList []*ArxivGroup) error {
	query := "SELECT id, original_arxiv_group_name FROM " + arxivGroupsTable
	var fetchedGroupList []*ArxivGroup
	err := pgxscan.Select(context.Background(), dbConnection.Pool, &fetchedGroupList, query)
	if err != nil {
		return fmt.Errorf("scanning the arXiv's group list: %w", err)
	}

	arxivGroupIdMapByName := make(map[string]ID)
	for _, arxivGroup := range fetchedGroupList {
		arxivGroupIdMapByName[arxivGroup.OriginalArxivGroupName] = arxivGroup.Id
	}

	for _, group := range groupList {
		if group.Id == 0 {
			group.Id = arxivGroupIdMapByName[group.OriginalArxivGroupName]
		}
		updateGroupReferenceInArxivArchive(group)
	}

	return nil
}

func updateGroupReferenceInArxivArchive(group *ArxivGroup) {
	for _, archive := range group.ArxivArchives {
		archive.ArxivGroupId = group.Id
	}
}
