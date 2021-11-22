package database

type Organisation struct {
	Id   ID     `db:"id"`
	Name string `db:"name"`
}
