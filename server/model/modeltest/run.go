package modeltest

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/test"
	"github.com/stretchr/testify/assert"
)

var resetTime = time.Date(2022, 06, 07, 13, 03, 24, 100, time.UTC)

func AssertRunEqual(t *testing.T, expected, actual test.Run) bool {
	t.Helper()

	// assert.Equal doesn't work on time.Time vars (see https://stackoverflow.com/a/69362528)
	// We could manually compare each field, but if we add a new field to the struct but not the test,
	// the test would still pass, even if the json encoding is incorrect.
	// So instead, we can reset all date fields and compare them separately.
	// If we add a new date field, the `assert.Equal(t, run, actual)` will catch it
	expectedCreatedAt := expected.CreatedAt
	expectedServiceTriggeredAt := expected.ServiceTriggeredAt
	expectedServiceTriggerCompletedAt := expected.ServiceTriggerCompletedAt
	expectedObtainedTraceAt := expected.ObtainedTraceAt
	expectedCompletedAt := expected.CompletedAt

	expected.CreatedAt = resetTime
	expected.ServiceTriggeredAt = resetTime
	expected.ServiceTriggerCompletedAt = resetTime
	expected.ObtainedTraceAt = resetTime
	expected.CompletedAt = resetTime
	if expected.Trace != nil {
		expected.Trace.RootSpan.StartTime = resetTime
		expected.Trace.RootSpan.EndTime = resetTime
		for _, span := range expected.Trace.Flat {
			span.StartTime = resetTime
			span.EndTime = resetTime
		}
	}

	actualCreatedAt := actual.CreatedAt
	actualServiceTriggeredAt := actual.ServiceTriggeredAt
	actualServiceTriggerCompletedAt := actual.ServiceTriggerCompletedAt
	actualObtainedTraceAt := actual.ObtainedTraceAt
	actualCompletedAt := actual.CompletedAt

	actual.CreatedAt = resetTime
	actual.ServiceTriggeredAt = resetTime
	actual.ServiceTriggerCompletedAt = resetTime
	actual.ObtainedTraceAt = resetTime
	actual.CompletedAt = resetTime
	if actual.Trace != nil {
		actual.Trace.RootSpan.StartTime = resetTime
		actual.Trace.RootSpan.EndTime = resetTime
		for _, span := range actual.Trace.Flat {
			span.StartTime = resetTime
			span.EndTime = resetTime
		}
	}

	return assert.Equal(t, expected, actual) &&

		assert.WithinDuration(t, expectedCreatedAt, actualCreatedAt, 0) &&
		assert.WithinDuration(t, expectedServiceTriggeredAt, actualServiceTriggeredAt, 0) &&
		assert.WithinDuration(t, expectedServiceTriggerCompletedAt, actualServiceTriggerCompletedAt, 0) &&
		assert.WithinDuration(t, expectedObtainedTraceAt, actualObtainedTraceAt, 0) &&
		assert.WithinDuration(t, expectedCompletedAt, actualCompletedAt, 0)
}
