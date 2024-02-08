package models

import "time"

type TestRun struct {
	TestID   string
	RunID    string
	Name     string
	Type     string
	Endpoint string
	Status   string
	When     time.Duration
}
