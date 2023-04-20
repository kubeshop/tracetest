package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"go.opentelemetry.io/otel/trace"
)

func (ro RunOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name     string
		Value    string
		SpanID   string
		Resolved bool
		Error    string
	}{
		Name:     ro.Name,
		Value:    ro.Value,
		SpanID:   ro.SpanID,
		Resolved: ro.Resolved,
		Error:    errToString(ro.Error),
	})
}

func (ro *RunOutput) UnmarshalJSON(data []byte) error {
	aux := struct {
		Name     string
		Value    string
		SpanID   string
		Resolved bool
		Error    string
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ro.Name = aux.Name
	ro.Value = aux.Value
	ro.SpanID = aux.SpanID
	ro.Resolved = aux.Resolved
	if err := stringToErr(aux.Error); err != nil {
		ro.Error = err
	}

	return nil
}

func (sar SpanAssertionResult) MarshalJSON() ([]byte, error) {
	sid := ""
	if sar.SpanID != nil {
		sid = sar.SpanID.String()
	}
	return json.Marshal(&struct {
		SpanID        *string
		ObservedValue string
		CompareErr    string
	}{
		SpanID:        &sid,
		ObservedValue: sar.ObservedValue,
		CompareErr:    errToString(sar.CompareErr),
	})
}

func (sar *SpanAssertionResult) UnmarshalJSON(data []byte) error {
	aux := struct {
		SpanID        string
		ObservedValue string
		CompareErr    string
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var sid *trace.SpanID
	if aux.SpanID != "" {
		s, err := trace.SpanIDFromHex(aux.SpanID)
		if err != nil {
			return err
		}
		sid = &s
	}

	sar.SpanID = sid
	sar.ObservedValue = aux.ObservedValue
	if err := stringToErr(aux.CompareErr); err != nil {
		if err.Error() == comparator.ErrNoMatch.Error() {
			err = comparator.ErrNoMatch
		}

		sar.CompareErr = err
	}

	return nil
}

func (na *NamedAssertions) UnmarshalJSON(data []byte) error {
	type encodedNamedAssertions struct {
		Name       string
		Assertions []Assertion
	}

	var namedAssertion encodedNamedAssertions
	if err := json.Unmarshal(data, &namedAssertion); err != nil {
		// this might be an []Assertion from the older format
		// try to parse []Assertions instead

		err = json.Unmarshal(data, &namedAssertion.Assertions)
		if err != nil {
			return err
		}
	}

	na.Name = namedAssertion.Name
	na.Assertions = namedAssertion.Assertions

	return nil
}

type encodedRun struct {
	ID                        int
	ShortID                   string
	TraceID                   string
	SpanID                    string
	State                     string
	LastErrorString           string
	CreatedAt                 time.Time
	ServiceTriggeredAt        time.Time
	ServiceTriggerCompletedAt time.Time
	ObtainedTraceAt           time.Time
	CompletedAt               time.Time
	TriggerResult             TriggerResult
	Trace                     *Trace
	Results                   *RunResults
	TestVersion               int
	Metadata                  map[string]string
}

func (r Run) MarshalJSON() ([]byte, error) {
	return json.Marshal(&encodedRun{
		ID:                        r.ID,
		TraceID:                   r.TraceID.String(),
		SpanID:                    r.SpanID.String(),
		State:                     string(r.State),
		LastErrorString:           errToString(r.LastError),
		CreatedAt:                 r.CreatedAt,
		ServiceTriggeredAt:        r.ServiceTriggeredAt,
		ServiceTriggerCompletedAt: r.ServiceTriggerCompletedAt,
		ObtainedTraceAt:           r.ObtainedTraceAt,
		CompletedAt:               r.CompletedAt,
		TestVersion:               r.TestVersion,
		Trace:                     r.Trace,
		Results:                   r.Results,
		TriggerResult:             r.TriggerResult,
		Metadata:                  r.Metadata,
	})
}

func (r *Run) UnmarshalJSON(data []byte) error {
	aux := encodedRun{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal run: %w", err)
	}

	var (
		tid trace.TraceID
		sid trace.SpanID
		err error
	)

	if aux.TraceID != "" {
		tid, err = trace.TraceIDFromHex(aux.TraceID)
		if err != nil {
			return fmt.Errorf("unmarshal run: %w", err)
		}
	}

	if aux.SpanID != "" {
		sid, err = trace.SpanIDFromHex(aux.SpanID)
		if err != nil {
			return fmt.Errorf("unmarshal run: %w", err)
		}
	}

	triggerResult := TriggerResult{
		Type: aux.TriggerResult.Type,
		HTTP: aux.TriggerResult.HTTP,
		GRPC: aux.TriggerResult.GRPC,
	}

	r.ID = aux.ID
	r.TraceID = tid
	r.SpanID = sid
	r.State = RunState(aux.State)
	r.LastError = stringToErr(aux.LastErrorString)
	r.CreatedAt = aux.CreatedAt
	r.ServiceTriggeredAt = aux.ServiceTriggeredAt
	r.ServiceTriggerCompletedAt = aux.ServiceTriggerCompletedAt
	r.ObtainedTraceAt = aux.ObtainedTraceAt
	r.CompletedAt = aux.CompletedAt
	r.TestVersion = aux.TestVersion
	r.TriggerResult = triggerResult

	r.Trace = aux.Trace
	r.Results = aux.Results
	r.Metadata = aux.Metadata

	return nil
}

func errToString(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}

func stringToErr(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}
