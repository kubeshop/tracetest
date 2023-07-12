package helpers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type tableColumnMetadata struct {
	Name string
	Len  int
}

func UnmarshalTable(t *testing.T, data string) []map[string]string {
	lines := strings.Split(data, "\n")
	require.GreaterOrEqual(t, len(lines), 2) // it should have at least two lines

	columnsMetadata := inferTableColumns(lines)
	result := []map[string]string{}

	lastLineIndexToParse := len(lines) - 2 // minus two, one to convert to index and another because the last line is always an empty line
	for i := 2; i <= lastLineIndexToParse; i++ {
		parsedResult := parseTableLine(columnsMetadata, lines[i])
		result = append(result, parsedResult)
	}

	return result
}

func inferTableColumns(tableLines []string) []tableColumnMetadata {
	headerLine := tableLines[0]
	separatorLine := tableLines[1]

	columns := strings.Split(separatorLine, " ")

	result := []tableColumnMetadata{}

	parserStart := 0
	for _, column := range columns {
		columnLen := len(column)

		columnName := strings.TrimSpace(headerLine[parserStart : parserStart+columnLen])
		parserStart = parserStart + columnLen + 1 // we need to shift one item for the next start

		result = append(result, tableColumnMetadata{
			Len:  columnLen,
			Name: columnName,
		})
	}

	return result
}

func parseTableLine(columnsMetadata []tableColumnMetadata, line string) map[string]string {
	result := map[string]string{}

	parserStart := 0
	for _, column := range columnsMetadata {
		value := strings.TrimSpace(line[parserStart : parserStart+column.Len])
		parserStart = parserStart + column.Len + 1

		result[column.Name] = value
	}

	return result
}
