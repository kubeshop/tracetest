package model_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
)

func TestRunExecutionTime(t *testing.T) {
	cases := []struct {
		name     string
		run      model.Run
		now      time.Time
		expected int
	}{
		{
			name: "CompletedOk",
			run: model.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, 100, time.UTC),
				CompletedAt: time.Date(2022, 01, 25, 12, 45, 36, 400, time.UTC),
			},
			expected: 4,
		},
		{
			name: "LessThan1Sec",
			run: model.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, 100, time.UTC),
				CompletedAt: time.Date(2022, 01, 25, 12, 45, 33, 400, time.UTC),
			},
			expected: 1,
		},
		{
			name: "StillRunning",
			run: model.Run{
				CreatedAt: time.Date(2022, 01, 25, 12, 45, 33, 100, time.UTC),
			},
			now:      time.Date(2022, 01, 25, 12, 45, 34, 300, time.UTC),
			expected: 2,
		},
		{
			name: "ZeroedDate",
			run: model.Run{
				CreatedAt:   time.Date(2022, 01, 25, 12, 45, 33, 100, time.UTC),
				CompletedAt: time.Unix(0, 0),
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

			assert.Equal(t, c.expected, c.run.ExecutionTime())
			model.Now = now
		})
	}
}
