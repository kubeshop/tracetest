package tracedb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

type sumologicDB struct {
	realTraceDB

	URL       string
	AccessID  string
	AccessKey string
}

type sumologicSpanSummary struct {
	ID        string `json:"id"`
	Name      string `json:"operationName"`
	ParentID  string `json:"parentId"`
	StartedAt string `json:"startedAt"`
	Duration  int64  `json:"duration"`
}

type getTraceSpansResponse struct {
	Page       []sumologicSpanSummary `json:"spanPage"`
	TotalCount int                    `json:"totalCount"`
	Next       string                 `json:"next"`
}

func NewSumoLogicDB(config *datastore.SumoLogicConfig) (TraceDB, error) {
	if config == nil {
		return nil, fmt.Errorf("empty config")
	}

	return &sumologicDB{
		URL:       config.URL,
		AccessID:  config.AccessID,
		AccessKey: config.AccessKey,
	}, nil
}

// Close implements TraceDB.
func (db *sumologicDB) Close() error {
	return nil
}

// Connect implements TraceDB.
func (db *sumologicDB) Connect(ctx context.Context) error {
	return nil
}

// GetEndpoints implements TraceDB.
func (db *sumologicDB) GetEndpoints() string {
	return db.URL
}

// GetTraceByID implements TraceDB.
func (db *sumologicDB) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	summaries, err := db.getTraceSpans(ctx, traceID, "")
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not get list of spans from trace: %w", err)
	}

	spans := db.convertSumoLogicSpanSummariesIntoSpans(summaries)
	return traces.NewTrace(traceID, spans), nil
}

func (db *sumologicDB) getTraceSpans(ctx context.Context, traceID string, token string) ([]sumologicSpanSummary, error) {
	spans := make([]sumologicSpanSummary, 0)
	response, err := db.getSpansPage(ctx, traceID, "")
	if err != nil {
		return nil, err
	}

	spans = append(spans, response.Page...)

	for response.Next != "" {
		response, err = db.getSpansPage(ctx, traceID, response.Next)
		if err != nil {
			return spans, err
		}

		spans = append(spans, response.Page...)
	}

	return spans, nil
}

func (db *sumologicDB) getSpansPage(ctx context.Context, traceID string, token string) (*getTraceSpansResponse, error) {
	url := fmt.Sprintf("/api/v1/tracing/traces/%s/spans?limit=100", traceID)
	if token != "" {
		url = fmt.Sprintf("%s&token=%s", url, token)
	}

	req, err := db.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute getTraceRequest: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, connection.ErrTraceNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code. Expected 200, got %d: %w", resp.StatusCode, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read getTraceSpans response body: %w", err)
	}

	var getTraceSpansResponse getTraceSpansResponse
	err = json.Unmarshal(body, &getTraceSpansResponse)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal getTraceSpans response body into struct: %w", err)
	}

	return &getTraceSpansResponse, nil
}

func (db *sumologicDB) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", db.URL, path), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create getTraceRequest: %w", err)
	}

	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", db.AccessID, db.AccessKey)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicAuth))

	return req, nil
}

func (db *sumologicDB) convertSumoLogicSpanSummariesIntoSpans(summaries []sumologicSpanSummary) []traces.Span {
	spans := make([]traces.Span, 0, len(summaries))
	for _, summary := range summaries {
		spanID, _ := trace.SpanIDFromHex(summary.ID)
		startTime, _ := time.Parse(time.RFC3339Nano, summary.StartedAt)
		endTime := startTime.Add(time.Duration(summary.Duration) * time.Nanosecond)

		spans = append(spans, traces.Span{
			ID:   spanID,
			Name: summary.Name,
			Attributes: traces.NewAttributes(map[string]string{
				"parent_id": summary.ParentID,
			}),
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	return spans
}

// Ready implements TraceDB.
func (db *sumologicDB) Ready() bool {
	return true
}

// AugmentTrace implements TraceAugmenter.
func (db *sumologicDB) AugmentTrace(ctx context.Context, trace *traces.Trace) (*traces.Trace, error) {
	if trace == nil {
		return nil, nil
	}

	spans := make([]traces.Span, 0, len(trace.Flat))
	for id, span := range trace.Flat {
		if span.Name == traces.TemporaryRootSpanName || span.Name == traces.TriggerSpanName {
			spans = append(spans, *span)
			continue
		}

		span, err := db.getAugmentedSpan(ctx, trace.ID.String(), id.String())
		if err != nil {
			return nil, err
		}

		spans = append(spans, *span)
	}

	newTrace := traces.NewTrace(trace.ID.String(), spans)

	return &newTrace, nil
}

type augmentedSpan struct {
	ID         string                `json:"id"`
	Name       string                `json:"operationName"`
	ParentID   string                `json:"parentId"`
	StartedAt  string                `json:"startedAt"`
	Duration   int64                 `json:"duration"`
	Attributes map[string]typedValue `json:"fields"`
	Events     []augmentedSpanEvent  `json:"events"`
}

type typedValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type augmentedSpanEvent struct {
	Timestamp  string           `json:"timestamp"`
	Name       string           `json:"name"`
	Attributes []eventAttribute `json:"attributes"`
}

type eventAttribute struct {
	Name  string     `json:"attributeName"`
	Value typedValue `json:"attributeValue"`
}

func (db *sumologicDB) getAugmentedSpan(ctx context.Context, traceID string, spanID string) (*traces.Span, error) {
	req, err := db.newRequest(http.MethodGet, fmt.Sprintf("/api/v1/tracing/traces/%s/spans/%s", traceID, spanID), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute augmented span: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		// We exceeded the rate limit, wait a bit and retry
		time.Sleep(10 * time.Second)
		return db.getAugmentedSpan(ctx, traceID, spanID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code. Expected 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	var span augmentedSpan
	err = json.Unmarshal(body, &span)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal augmented span into struct: %w", err)
	}

	id, _ := trace.SpanIDFromHex(span.ID)
	startTime, _ := time.Parse(time.RFC3339Nano, span.StartedAt)
	endTime := startTime.Add(time.Duration(span.Duration) * time.Nanosecond)

	attributes := map[string]string{
		"parent_id": span.ParentID,
	}
	for name, typedValue := range span.Attributes {
		attributes[name] = typedValue.Value
	}

	events := make([]traces.SpanEvent, 0, len(span.Events))
	for _, event := range span.Events {
		timestamp, _ := time.Parse(time.RFC3339Nano, event.Timestamp)
		eventAttributes := make(map[string]string, len(event.Attributes))
		for _, attribute := range event.Attributes {
			eventAttributes[attribute.Name] = attribute.Value.Value
		}

		events = append(events, traces.SpanEvent{
			Timestamp:  timestamp,
			Name:       event.Name,
			Attributes: traces.NewAttributes(eventAttributes),
		})
	}

	return &traces.Span{
		ID:         id,
		Name:       span.Name,
		StartTime:  startTime,
		EndTime:    endTime,
		Attributes: traces.NewAttributes(attributes),
		Events:     events,
	}, nil
}

var _ TraceDB = &sumologicDB{}
var _ TraceAugmenter = &sumologicDB{}
