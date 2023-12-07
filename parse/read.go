package parse

import (
	"fmt"
	"io"

	"github.com/xwb1989/sqlparser"
)

func Read(r io.Reader) *Schema {
	schema := &Schema{
		TableMap: map[string]Table{},
		Tables:   []string{},
	}
	tokens := sqlparser.NewTokenizer(r)

	for {
		stmt, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err.Error())
		}

		switch v := stmt.(type) {
		case *sqlparser.DDL:
			schema.ddl(v)
		}
	}

	return schema
}

func (s *Schema) ddl(v *sqlparser.DDL) {
	switch v.Action {
	case sqlparser.CreateStr:
		s.table(v)
	}
}

func (s *Schema) table(v *sqlparser.DDL) {
	tableName := v.NewName.Name.String()
	s.Tables = append(s.Tables, tableName)
	t := Table{
		Name:      tableName,
		ColumnMap: map[string]Column{},
		Columns:   []string{},
	}

	if v.TableSpec != nil {
		for _, colSpec := range v.TableSpec.Columns {
			columnName := colSpec.Name.String()
			c := Column{
				Name: columnName,
				Type: colSpec.Type.Type,
			}

			t.Columns = append(t.Columns, columnName)
			t.ColumnMap[columnName] = c
		}
	}

	s.TableMap[tableName] = t
}
