package comparator_test

import (
	"testing"

	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {
	assert.Equal(t, "=", comparator.Eq.String())

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, comparator.Eq.Compare("a", "a"))
		assert.ErrorIs(t, comparator.Eq.Compare("a", "b"), comparator.ErrNoMatch)
	})

	t.Run("Int", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, comparator.Eq.Compare("1", "1"))
		assert.ErrorIs(t, comparator.Eq.Compare("", "2"), comparator.ErrNoMatch)
	})
}

func TestGt(t *testing.T) {
	assert.Equal(t, ">", comparator.Gt.String())

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, comparator.Gt.Compare("a", "1"), `cannot parse "a" as integer`)
	})

	t.Run("Int", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, comparator.Gt.Compare("2", "1"))
		assert.NoError(t, comparator.Gt.Compare("10", "2"))
		assert.ErrorIs(t, comparator.Gt.Compare("1", "1"), comparator.ErrNoMatch)
		assert.ErrorIs(t, comparator.Gt.Compare("1", "2"), comparator.ErrNoMatch)
	})
}

func TestLt(t *testing.T) {
	assert.Equal(t, "<", comparator.Lt.String())

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, comparator.Lt.Compare("a", "1"), `cannot parse "a" as integer`)
	})

	t.Run("Int", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, comparator.Lt.Compare("1", "2"))
		assert.NoError(t, comparator.Lt.Compare("2", "10"))
		assert.ErrorIs(t, comparator.Lt.Compare("1", "1"), comparator.ErrNoMatch)
		assert.ErrorIs(t, comparator.Lt.Compare("2", "1"), comparator.ErrNoMatch)
	})
}

func TestContains(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "contains", comparator.Contains.String())

	assert.NoError(t, comparator.Contains.Compare("lo", "hello"))
	assert.ErrorIs(t, comparator.Contains.Compare("hello", "lo"), comparator.ErrNoMatch)
}
