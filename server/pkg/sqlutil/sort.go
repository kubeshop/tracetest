package sqlutil

import (
	"fmt"
	"strings"
)

func Sort(sql, sortBy, sortDirection, defaultSortBy string, sortingFields map[string]string) string {
	sortField, ok := sortingFields[sortBy]

	if !ok {
		sortField = sortingFields[defaultSortBy]
	}

	dir := "DESC"
	if strings.ToLower(sortDirection) == "asc" {
		dir = "ASC"
	}

	return fmt.Sprintf("%s ORDER BY %s %s", sql, sortField, dir)
}
