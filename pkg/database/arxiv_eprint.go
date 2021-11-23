package database

import "time"

type ArxivEprint struct {
	Id ID `db:"id"`

	ArxivId string `db:"arxiv_id"`

	Comment       *string                 `db:"comment"`
	Extra         *map[string]interface{} `db:"extra"`
	LatestVersion int                     `db:"latest_version"`
	PdfLink       *string                 `db:"pdf_link"`

	PublishedAt time.Time `db:"published_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	EprintId ID `db:"eprint_id"`

	EPrint               *Eprint
	PrimaryArxivCategory *ArxivCategory
	OtherArxivCategories []*ArxivCategory
}
