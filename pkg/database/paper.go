package database

type Paper struct {
	Id ID `db:"id"`

	Doi        string `db:"doi"`
	JournalRef string `db:"journal_ref"`

	Abstract string `db:"abstract"`
	Title    string `db:"title"`
	Year     int    `db:"year"`

	Authors []*Author
}
