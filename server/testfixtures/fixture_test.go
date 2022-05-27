package testfixtures_test

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/testfixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixture(t *testing.T) {
	testfixtures.RegisterFixture("uuid", func(options testfixtures.FixtureOptions) (string, error) {
		return uuid.NewString(), nil
	})

	var uuid1, uuid2 string
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		uuid, err := testfixtures.GetFixtureValue[string]("uuid")
		require.NoError(t, err)
		uuid1 = uuid
		wg.Done()
	}()

	go func() {
		uuid, err := testfixtures.GetFixtureValue[string]("uuid")
		require.NoError(t, err)
		uuid2 = uuid
		wg.Done()
	}()

	wg.Wait()

	assert.Equal(t, uuid1, uuid2)
}
