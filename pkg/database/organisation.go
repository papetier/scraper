package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Organisation struct {
	Id   ID     `db:"id"`
	Name string `db:"name"`
}

const organisationsTable = "organisations"

var organisationsColumns = []string{
	"id",
	"name",
}

func saveOrganisationsTx(tx pgx.Tx, organisationList []*Organisation) error {
	log.Debug("saving organisations")

	var organisationValues []interface{}
	for _, organisation := range organisationList {
		organisationValues = append(organisationValues, organisation.Name)
	}

	organisationPlaceholder := generateInsertPlaceholder(len(organisationsColumns[1:]), len(organisationList), 1)
	organisationsQuery := "INSERT INTO " + organisationsTable + " (" + strings.Join(organisationsColumns[1:], ", ") + ") VALUES " + organisationPlaceholder + " ON CONFLICT DO NOTHING RETURNING id"

	organisationRows, err := tx.Query(context.Background(), organisationsQuery, organisationValues...)
	defer organisationRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the organisations into the database: %w", err)
	}

	insertedOrganisationCount := 0
	var insertedOrganisationIdList []ID
	for organisationRows.Next() {
		var id ID
		err = organisationRows.Scan(&id)
		insertedOrganisationIdList = append(insertedOrganisationIdList, id)
		if err != nil {
			return fmt.Errorf("scanning the organisation ids: %w", err)
		}
		insertedOrganisationCount++
	}

	if insertedOrganisationCount == len(organisationList) {
		for i, id := range insertedOrganisationIdList {
			organisationList[i].Id = id
		}
	} else {
		err = fetchAndUpdateOrganisationIdsTx(tx, organisationList)
		if err != nil {
			return fmt.Errorf("fetching the organisation ids: %w", err)
		}
	}

	return nil
}

func fetchAndUpdateOrganisationIdsTx(tx pgx.Tx, organisationList []*Organisation) error {
	placeholder := generateInsertPlaceholder(len(organisationList), 1, 1)
	query := "SELECT id, full_name FROM " + authorsTable + " WHERE full_name IN " + placeholder
	var parameters []interface{}
	for _, organisation := range organisationList {
		parameters = append(parameters, organisation.Name)
	}
	var fetchedOrganisationList []*Organisation
	err := pgxscan.Select(context.Background(), tx, &fetchedOrganisationList, query, parameters...)
	if err != nil {
		return fmt.Errorf("scanning the organisation list: %w", err)
	}

	organisationIdMapByName := make(map[string]ID)
	for _, organisation := range fetchedOrganisationList {
		organisationIdMapByName[organisation.Name] = organisation.Id
	}

	for _, organisation := range organisationList {
		if organisation.Id == 0 {
			organisation.Id = organisationIdMapByName[organisation.Name]
		}
	}

	return nil
}
