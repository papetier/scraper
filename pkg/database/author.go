package database

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Author struct {
	Id       ID      `db:"id"`
	Email    *string `db:"email"`
	FullName string  `db:"full_name"`

	Organisations []*Organisation
}

const (
	authorsTable              = "authors"
	authorsOrganisationsTable = "authors_organisations"
)

var authorsColumns = []string{
	"id",
	"email",
	"full_name",
}

var authorsOrganisationsColumns = []string{
	"author_id",
	"organisation_id",
}

func saveAuthorsWithOrganisationsTx(tx pgx.Tx, authorList []*Author) error {
	log.Debug("saving authors with their organisations")

	// get unique organisation list
	var organisationList []*Organisation
	organisationSet := make(map[string]struct{})
	for _, author := range authorList {
		for _, organisation := range author.Organisations {
			if _, exists := organisationSet[organisation.Name]; !exists {
				organisationList = append(organisationList, author.Organisations...)
				organisationSet[organisation.Name] = struct{}{}
			}
		}
	}

	// save all authors' organisations
	if len(organisationList) > 0 {
		err := saveOrganisationsTx(tx, organisationList)
		if err != nil {
			return fmt.Errorf("saving the author's organisations: %w", err)
		}
	}

	// re-map author's organisations
	// if several authors had same organisation
	for _, author := range authorList {
		for i, authorOrganisation := range author.Organisations {
			for _, organisation := range organisationList {
				if organisation.Name == authorOrganisation.Name {
					author.Organisations[i] = organisation
					break
				}
			}
		}
	}

	// save authors
	err := saveAuthorsTx(tx, authorList)
	if err != nil {
		return fmt.Errorf("saving the authors: %w", err)
	}

	// save the authors/organisations links
	if len(organisationList) > 0 {
		err = saveAuthorsOrganisationsTx(tx, authorList)
		if err != nil {
			return fmt.Errorf("saving the authors_organisations links: %w", err)
		}
	}

	return nil
}

func saveAuthorsTx(tx pgx.Tx, authorList []*Author) error {
	log.Debug("saving authors")

	var authorValues []interface{}
	for _, author := range authorList {
		authorValues = append(authorValues, author.Email, author.FullName)
	}

	authorPlaceholder := generateInsertPlaceholder(len(authorsColumns[1:]), len(authorList), 1)
	authorsQuery := "INSERT INTO " + authorsTable + " (" + strings.Join(authorsColumns[1:], ", ") + ") VALUES " + authorPlaceholder + " ON CONFLICT DO NOTHING RETURNING id"

	authorRows, err := tx.Query(context.Background(), authorsQuery, authorValues...)
	defer authorRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the authors into the database: %w", err)
	}

	insertedAuthorCount := 0
	var insertedAuthorIdList []ID
	for authorRows.Next() {
		var id ID
		err = authorRows.Scan(&id)
		insertedAuthorIdList = append(insertedAuthorIdList, id)
		if err != nil {
			return fmt.Errorf("scanning the author ids: %w", err)
		}
		insertedAuthorCount++
	}

	if insertedAuthorCount == len(authorList) {
		for i, id := range insertedAuthorIdList {
			authorList[i].Id = id
		}
	} else {
		err = fetchAndUpdateAuthorIdsTx(tx, authorList)
		if err != nil {
			return fmt.Errorf("fetching the author ids: %w", err)
		}
	}

	return nil
}

func fetchAndUpdateAuthorIdsTx(tx pgx.Tx, authorList []*Author) error {
	placeholder := generateInsertPlaceholder(len(authorList), 1, 1)
	query := "SELECT id, full_name FROM " + authorsTable + " WHERE full_name IN " + placeholder
	var parameters []interface{}
	for _, author := range authorList {
		parameters = append(parameters, author.FullName)
	}
	var fetchedAuthorList []*Author
	err := pgxscan.Select(context.Background(), tx, &fetchedAuthorList, query, parameters...)
	if err != nil {
		return fmt.Errorf("scanning the author list: %w", err)
	}

	authorIdMapByName := make(map[string]ID)
	for _, author := range fetchedAuthorList {
		authorIdMapByName[author.FullName] = author.Id
	}

	for _, author := range authorList {
		if author.Id == 0 {
			author.Id = authorIdMapByName[author.FullName]
		}
	}

	return nil
}

func saveAuthorsOrganisationsTx(tx pgx.Tx, authorList []*Author) error {
	log.Debug("saving the authors_organisations links")

	linkCount := 0
	var authorOrganisationLinkValues []interface{}
	for _, author := range authorList {
		for _, organisation := range author.Organisations {
			linkCount++
			authorOrganisationLinkValues = append(authorOrganisationLinkValues, author.Id, organisation.Id)
		}
	}

	authorLinkPlaceholder := generateInsertPlaceholder(len(authorsOrganisationsColumns), linkCount, 1)
	authorLinksQuery := "INSERT INTO " + authorsOrganisationsTable + " (" + strings.Join(authorsOrganisationsColumns, ", ") + ") VALUES " + authorLinkPlaceholder + " ON CONFLICT DO NOTHING"

	authorLinkRows, err := tx.Query(context.Background(), authorLinksQuery, authorOrganisationLinkValues...)
	defer authorLinkRows.Close()
	if err != nil {
		return fmt.Errorf("inserting the authors_organisations links into the database: %w", err)
	}

	return nil
}
