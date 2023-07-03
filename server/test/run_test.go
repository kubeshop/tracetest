package test_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/test"
	"github.com/stretchr/testify/assert"
)

func TestRunExecutionTime(t *testing.T) {
	cases := []struct {
		name     string
		run      test.Run
		now      time.Time
		expected int
	}{
		{
			name: "CompletedOk",
			run: test.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Date(2022, 01, 25, 12, 45, 36, int(400*time.Millisecond), time.UTC),
			},
			expected: 4,
		},
		{
			name: "LessThan1Sec",
			run: test.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Date(2022, 01, 25, 12, 45, 33, int(400*time.Millisecond), time.UTC),
			},
			expected: 1,
		},
		{
			name: "StillRunning",
			run: test.Run{
				CreatedAt: time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 2,
		},
		{
			name: "ZeroedDate",
			run: test.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				CompletedAt: time.Unix(0, 0),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := test.Now
			if c.now.Unix() > 0 {
				test.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, c.run.ExecutionTime())
			test.Now = now
		})
	}
}

func TestRunTriggerTime(t *testing.T) {
	cases := []struct {
		name     string
		run      test.Run
		now      time.Time
		expected int
	}{
		{
			name: "CompletedOk",
			run: test.Run{
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Date(2022, 01, 25, 12, 45, 36, int(400*time.Millisecond), time.UTC),
			},
			expected: 3300,
		},
		{
			name: "LessThan1Sec",
			run: test.Run{
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Date(2022, 01, 25, 12, 45, 33, int(400*time.Millisecond), time.UTC),
			},
			expected: 300,
		},
		{
			name: "StillRunning",
			run: test.Run{
				ServiceTriggeredAt: time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 1200,
		},
		{
			name: "ZeroedDate",
			run: test.Run{
				ServiceTriggeredAt:        time.Date(2022, 01, 25, 12, 45, 33, int(100*time.Millisecond), time.UTC),
				ServiceTriggerCompletedAt: time.Unix(0, 0),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, int(300*time.Millisecond), time.UTC),
			expected: 1200,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now := test.Now
			if c.now.Unix() > 0 {
				test.Now = func() time.Time {
					return c.now
				}
			}

			assert.Equal(t, c.expected, c.run.TriggerTime())
			test.Now = now
		})
	}
}
