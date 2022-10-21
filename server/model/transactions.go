package model

import (
	"time"

	"github.com/kubeshop/tracetest/server/id"
)

type (
	Transaction struct {
		ID          id.ID
		CreatedAt   time.Time
		Name        string
		Description string
		Version     int
		Steps       []Test
		Summary     Summary
	}
)

func (t Transaction) HasID() bool {
	return t.ID != ""
}
