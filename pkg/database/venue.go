package database

type Venue struct {
	Id   ID     `db:"id"`
	Name string `db:"name"`
}
