package wizard

import "time"

type Wizard struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	ID          string     `json:"id"`
	State       StepState  `json:"state"`
	CompletedAt *time.Time `json:"completed_at"`
}

type StepState string

const (
	StepStatusPending StepState = "pending"
	StepStatusRunning StepState = "inProgress"
	StepStatusDone    StepState = "completed"
)
