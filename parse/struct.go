package parse

type Schema struct {
	// Map table name to table struct
	TableMap map[string]Table

	// List of table names
	Tables []string
}

type Table struct {
	Name string

	// Map column name to column struct
	ColumnMap map[string]Column

	// List of table columns
	Columns []string
}

type Column struct {
	Name string
	Type string
}
