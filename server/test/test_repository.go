package test

import (
	"database/sql"
)

type Repository interface {
	// TODO: uncomment as we add new functionality

	// Create(context.Context, Test) (Test, error)
	// Update(context.Context, Test) (Test, error)
	// Delete(context.Context, Test) error
	// Exists(context.Context, id.ID) (bool, error)
	// GetLatestTestVersion(context.Context, id.ID) (Test, error)
	// GetTestVersion(_ context.Context, _ id.ID, version int) (Test, error)
	// GetTests(_ context.Context, take, skip int32, query, sortBy, sortDirection string) ([]Test, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{db}
}
