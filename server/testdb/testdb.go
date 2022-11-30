package testdb

import (
	"errors"

	"github.com/kubeshop/tracetest/server/id"
)

var (
	IDGen       = id.NewRandGenerator()
	ErrNotFound = errors.New("record not found")
)
