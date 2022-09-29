package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

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

func (a Assertion) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Attribute  string
		Comparator string
		Value      *AssertionExpression
	}{
		Attribute:  a.Attribute.String(),
		Comparator: a.Comparator.String(),
		Value:      a.Value,
	})
}

func (a *Assertion) UnmarshalJSON(data []byte) error {
	aux := struct {
		Attribute  string
		Comparator string
		Value      *AssertionExpression
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		// This might be an older assertion that has its value as a string
		// we have to try convert it before returning the error
		oldAssertion := struct {
			Attribute  string
			Comparator string
			Value      string
		}{}

		if err := json.Unmarshal(data, &oldAssertion); err != nil {
			// It isn't the old format either, so it's an error
			return err
		}

		aux.Attribute = oldAssertion.Attribute
		aux.Comparator = oldAssertion.Comparator
		aux.Value = &AssertionExpression{
			LiteralValue: LiteralValue{
				Type:  "string",
				Value: oldAssertion.Value,
			},
		}
	}

	c, err := comparator.DefaultRegistry().Get(aux.Comparator)
	if err != nil {
		return err
	}

	a.Attribute = Attribute(aux.Attribute)
	a.Value = aux.Value
	a.Comparator = c

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
	Trace                     *traces.Trace
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
