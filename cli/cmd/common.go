package cmd

import (
	"fmt"
	"time"

	"github.com/Jeffail/gabs/v2"
)

func formatItemDate(item *gabs.Container, path string) error {
	rawDate := item.Path(path).Data()
	if rawDate == nil {
		return nil
	}
	dateStr := rawDate.(string)
	// if field is empty, do nothing
	if dateStr == "" {
		return nil
	}

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse datetime field '%s' (value '%s'): %s", path, dateStr, err)
	}

	if date.IsZero() {
		// sometime the date comes like 0000-00-00T00:00:00Z... show nothing in that case
		item.SetP("", path)
		return nil
	}

	item.SetP(date.Format(time.DateTime), path)
	return nil
}
