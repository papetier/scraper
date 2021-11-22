package database

type Eprint struct {
	Id      ID `db:"id"`
	PaperId ID `db:"paper_id"`

	Paper *Paper
}
