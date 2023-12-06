package tracedb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/stretchr/testify/require"
)

func TestSumoLogic(t *testing.T) {
	db := tracedb.NewSumoLogicDB()

	trace, err := db.GetTraceByID(context.Background(), "5c4a84746737e161f1ddcbfe3087255b")
	require.NoError(t, err)
	require.True(t, trace.RootSpan.ID.IsValid())
	require.Greater(t, len(trace.Flat), 0)

	augmenter := db.(tracedb.TraceAugmenter)
	newTrace, err := augmenter.AugmentTrace(context.Background(), &trace)
	require.NoError(t, err)
}
