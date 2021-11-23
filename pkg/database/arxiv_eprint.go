package database

type ArxivEprint struct {
	Id ID `db:"id"`

	ArxivId string `db:"arxiv_id"`

	Comment       string                 `db:"comment"`
	Extra         map[string]interface{} `db:"extra"`
	LatestVersion string                 `db:"latest_version"`
	PdfLink       string                 `db:"pdf_link"`

	PublishedAt string `db:"published_at"`
	UpdatedAt   string `db:"updated_at"`

	EprintId ID `db:"eprint_id"`

	EPrint          *Eprint
	ArxivCategories []*ArxivCategory
}
