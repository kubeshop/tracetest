package model

import (
	"encoding/json"
	"errors"
)

type ConnectionResult struct {
	PortCheck      ConnectionTestStep
	Connectivity   ConnectionTestStep
	Authentication ConnectionTestStep
	FetchTraces    ConnectionTestStep
}

func (c ConnectionResult) HasSucceed() bool {
	return c.Connectivity.HasSucceed() && c.Authentication.HasSucceed() && c.FetchTraces.HasSucceed()
}

type ConnectionTestStep struct {
	Passed  bool
	Status  Status
	Message string
	Error   error
}

func (s *ConnectionTestStep) UnmarshalJSON(bytes []byte) error {
	var step struct {
		Passed       bool
		Status       Status
		Message      string
		ErrorMessage string
	}

	err := json.Unmarshal(bytes, &step)
	if err != nil {
		return err
	}

	s.Passed = step.Passed
	s.Status = step.Status
	s.Message = step.Message
	if step.ErrorMessage != "" {
		s.Error = errors.New(step.ErrorMessage)
	}

	return nil
}

func (r *ConnectionTestStep) HasSucceed() bool {
	if r == nil {
		return true
	}

	return r.Error == nil
}

func (r *ConnectionTestStep) IsSet() bool {
	if r == nil {
		return false
	}

	return r.Message != ""
}
