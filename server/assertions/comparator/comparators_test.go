package comparator_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	t.Run("CannotHaveDuplicates", func(t *testing.T) {
		t.Parallel()
		f, err := comparator.NewRegistry(
			comparator.Eq,
			comparator.Eq,
		)
		assert.Nil(t, f)
		assert.EqualError(t, err, `comparator "=" already registered`)
	})

	t.Run("GetsCorrectComparator", func(t *testing.T) {
		t.Parallel()
		f, err := comparator.NewRegistry(comparator.Basic...)
		assert.NoError(t, err)

		c, err := f.Get(">")
		assert.NoError(t, err)
		assert.Equal(t, ">", c.String())

		c, err = f.Get("!=`")
		assert.Nil(t, c)
		assert.ErrorIs(t, err, comparator.ErrNotFound)
	})
}
