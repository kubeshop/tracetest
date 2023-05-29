package sqlutil

import "strings"

func Search(sql, condition, query string, params []any) (string, []any) {
	if query == "" {
		return sql, params
	}

	sql += condition
	params = append(params, CleanSearchQuery(query))

	return sql, params
}

func CleanSearchQuery(query string) string {
	return "%" + strings.ReplaceAll(query, " ", "%") + "%"
}
