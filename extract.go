package sqrl

func ExtractColumns(builder Sqlizer) []string {
	switch b := builder.(type) {
	case *UpdateBuilder:
		cols := make([]string, len(b.setClauses))
		for i, clause := range b.setClauses {
			cols[i] = clause.column
		}
		return cols
	case *InsertBuilder:
		return b.columns
	default:
		panic("failed to extract columns")
	}
}

func ExtractValues(builder Sqlizer) [][]interface{} {
	switch b := builder.(type) {
	case *UpdateBuilder:
		vals := make([]interface{}, len(b.setClauses))
		for i, clause := range b.setClauses {
			vals[i] = clause.value
		}
		return [][]interface{}{vals}
	case *InsertBuilder:
		return b.values
	default:
		panic("failed to extract values")
	}
}

func ExtractWhereParts(builder Sqlizer) []Sqlizer {
	switch b := builder.(type) {
	case *SelectBuilder:
		return b.whereParts
	case *UpdateBuilder:
		return b.whereParts
	case *DeleteBuilder:
		return b.whereParts
	case *InsertBuilder:
		return nil
	default:
		panic("failed to extract whereParts")
	}
}

func ExtractTableNames(builder Sqlizer) []string {
	switch b := builder.(type) {
	case *SelectBuilder:
		tables := make([]string, 0, len(b.fromParts))
		for _, from := range b.fromParts {
			switch from := from.(type) {
			case *part:
				tables = append(tables, from.pred.(string))
			case aliasExpr:
				tables = append(tables, from.alias)
			}
		}
		return tables
	case *InsertBuilder:
		return []string{b.into}
	case *UpdateBuilder:
		return []string{b.table}
	case *DeleteBuilder:
		return []string{b.from}
	default:
		return nil
	}
}
