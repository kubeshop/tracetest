package testdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerIDForFirstTime(t *testing.T) {
	db, clean := getDB()
	defer clean()

	id, isNew, err := db.ServerID()
	assert.Len(t, id, 10)
	assert.True(t, isNew)
	assert.NoError(t, err)
}

func TestGetServerIDForSecondTime(t *testing.T) {
	db, clean := getDB()
	defer clean()

	firstId, isNew, err := db.ServerID()
	assert.Len(t, firstId, 10)
	assert.True(t, isNew)
	assert.NoError(t, err)

	secondId, isNew, err := db.ServerID()
	assert.Len(t, secondId, 10)
	assert.False(t, isNew)
	assert.NoError(t, err)

	assert.Equal(t, firstId, secondId)
}
