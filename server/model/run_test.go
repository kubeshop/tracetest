package model_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
)

func TestRunExecutionTime(t *testing.T) {
	timeTest(t, func(start, end time.Time) int {
		run := model.Run{
			CreatedAt:   start,
			CompletedAt: end,
		}

		return run.ExecutionTime()
	})
}

func TestRunTriggerTime(t *testing.T) {
	timeTest(t, func(start, end time.Time) int {
		run := model.Run{
			ServiceTriggeredAt:        start,
			ServiceTriggerCompletedAt: end,
		}

		return run.TriggerTime()
	})
}

func timeTest(t *testing.T, fn func(start, end time.Time) int) {
	t.Helper()

	type timeCase struct {
		start, end time.Time
	}
	cases := []struct {
		name     string
		run      timeCase
		now      time.Time
		expected int
	}{
		{
			name: "CompletedOk",
			run: timeCase{
				start: time.Date(2022, 01, 25, 12, 45, 33, 100, time.UTC),
				end:   time.Date(2022, 01, 25, 12, 45, 36, 400, time.UTC),
			},
			expected: 4,
		},
		{
			name: "LessThan1Sec",
			run: timeCase{
				start: time.Date(2022, 01, 25, 12, 45, 33, 100, time.UTC),
				end:   time.Date(2022, 01, 25, 12, 45, 33, 400, time.UTC),
			},
			expected: 1,
		},
		{
			name: "StillRunning",
			run: timeCase{
				start: time.Date(2022, 01, 25, 12, 45, 33, 100, time.UTC),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, 300, time.UTC),
			expected: 2,
		},
		{
			name: "ZeroedDate",
			run: timeCase{
				start: time.Date(2022, 01, 25, 12, 45, 33, 100, time.UTC),
				end:   time.Unix(0, 0),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, 300, time.UTC),
			expected: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := model.Now
			if c.now.Unix() > 0 {
				model.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, fn(c.run.start, c.run.end))
			model.Now = now
		})
	}
}
