package tests_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/tests"
	"github.com/stretchr/testify/assert"
)

func TestRunExecutionTime(t *testing.T) {
	cases := []struct {
		name     string
		run      tests.Run
		now      time.Time
		expected int
	}{
		{
			name: "CompletedOk",
			run: tests.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Date(2022, 01, 25, 12, 45, 36, int(400*time.Millisecond), time.UTC),
			},
			expected: 4,
		},
		{
			name: "LessThan1Sec",
			run: tests.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Date(2022, 01, 25, 12, 45, 33, int(400*time.Millisecond), time.UTC),
			},
			expected: 1,
		},
		{
			name: "StillRunning",
			run: tests.Run{
				CreatedAt: time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 2,
		},
		{
			name: "ZeroedDate",
			run: tests.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Unix(0, 0),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := tests.Now
			if c.now.Unix() > 0 {
				tests.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, c.run.ExecutionTime())
			tests.Now = now
		})
	}
}

func TestRunTriggerTime(t *testing.T) {
	cases := []struct {
		name     string
		run      tests.Run
		now      time.Time
		expected int
	}{
		{
			name: "CompletedOk",
			run: tests.Run{
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Date(2022, 01, 25, 12, 45, 36, int(400*time.Millisecond), time.UTC),
			},
			expected: 3300,
		},
		{
			name: "LessThan1Sec",
			run: tests.Run{
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Date(2022, 01, 25, 12, 45, 33, int(400*time.Millisecond), time.UTC),
			},
			expected: 300,
		},
		{
			name: "StillRunning",
			run: tests.Run{
				ServiceTriggeredAt: time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 1200,
		},
		{
			name: "ZeroedDate",
			run: tests.Run{
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Unix(0, 0),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 1200,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := tests.Now
			if c.now.Unix() > 0 {
				tests.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, c.run.TriggerTime())
			tests.Now = now
		})
	}
}
