package types

import (
	"time"

	"github.com/kubeshop/tracetest/server/pkg/maps"
)

type TestResource struct {
	Type string `json:"type"`
	Spec Test   `json:"spec"`
}

type Test struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Version     int                          `json:"version"`
	Trigger     Trigger                      `json:"trigger"`
	Specs       []TestSpec                   `json:"specs"`
	Outputs     maps.Ordered[string, Output] `json:"outputs"`
	Summary     Summary                      `json:"summary"`
}

type Trigger struct {
	Type        string      `json:"type"`
	HTTPRequest HTTPRequest `json:"http"`
}

type HTTPRequest struct {
	Method  string       `json:"method,omitempty"`
	URL     string       `json:"url"`
	Body    string       `json:"body,omitempty"`
	Headers []HTTPHeader `json:"headers,omitempty"`
}

type HTTPHeader struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type Output struct {
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

type TestSpec struct {
	Selector   Selector `json:"selector"`
	Name       string   `json:"name,omitempty"`
	Assertions []string `json:"assertions"`
}

type Selector struct {
	Query          string       `json:"query"`
	ParsedSelector SpanSelector `json:"parsedSelector"`
}

type SpanSelector struct {
	Filters       []SelectorFilter     `json:"filters"`
	PseudoClass   *SelectorPseudoClass `json:"pseudoClass,omitempty"`
	ChildSelector *SpanSelector        `json:"childSelector,omitempty"`
}

type SelectorFilter struct {
	Property string `json:"property"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type SelectorPseudoClass struct {
	Name     string `json:"name"`
	Argument *int32 `json:"argument,omitempty"`
}

type Summary struct {
	Runs    int     `json:"runs"`
	LastRun LastRun `json:"lastRun"`
}

type LastRun struct {
	Time   time.Time `json:"time"`
	Passes int       `json:"passes"`
	Fails  int       `json:"fails"`
}
