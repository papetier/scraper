package database

type Publisher struct {
	Id   ID     `db:"id"`
	Name string `db:"name"`
	Url  string `db:"url"`
}
