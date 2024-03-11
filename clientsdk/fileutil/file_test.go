package fileutil_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetID(t *testing.T) {
	t.Run("NoID", func(t *testing.T) {
		f, err := fileutil.Read("../../testdata/definitions/valid_http_test_definition.yml")
		require.NoError(t, err)

		f, err = f.SetID("new-id")
		require.NoError(t, err)

		assert.True(t, f.HasID())
	})

	t.Run("WithID", func(t *testing.T) {
		f, err := fileutil.Read("../../testdata/definitions/valid_http_test_definition_with_id.yml")
		require.NoError(t, err)

		_, err = f.SetID("new-id")
		require.ErrorIs(t, err, fileutil.ErrFileHasID)
	})

}
