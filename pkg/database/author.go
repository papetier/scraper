package database

type Author struct {
	Id       ID     `db:"id"`
	Email    string `db:"email"`
	FullName string `db:"full_name"`

	Organisations []*Organisation
}
