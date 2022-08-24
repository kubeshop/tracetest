package tracedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

type signalfxDB struct {
	Token string
	Realm string
	URL   string

	httpClient *http.Client
}

func (db signalfxDB) getURL() string {
	if db.URL != "" {
		return db.URL
	}

	return fmt.Sprintf("https://api.%s.signalfx.com", db.Realm)
}

func (db signalfxDB) Close() error {
	// Doesn't need to be closed
	return nil
}

func (db signalfxDB) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	timestamps, err := db.getSegmentsTimestamps(ctx, traceID)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("coult not get trace segment timestamps: %w", err)
	}

	if len(timestamps) == 0 {
		return traces.Trace{}, ErrTraceNotFound
	}

	traceSpans := make([]traces.Span, 0)

	for _, timestamp := range timestamps {
		segmentSpans, err := db.getSegmentSpans(ctx, traceID, timestamp)
		if err != nil {
			return traces.Trace{}, fmt.Errorf("could not get segment spans: %w", err)
		}

		for _, signalFxSpan := range segmentSpans {
			span := convertSignalFXSpan(signalFxSpan)
			traceSpans = append(traceSpans, span)
		}
	}

	if len(traceSpans) == 0 {
		return traces.Trace{}, ErrTraceNotFound
	}

	return traces.New(traceID, traceSpans), nil
}

func (db signalfxDB) getSegmentsTimestamps(ctx context.Context, traceID string) ([]int64, error) {
	url := fmt.Sprintf("%s/v2/apm/trace/%s/segments", db.getURL(), traceID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []int64{}, fmt.Errorf("could not create request: %w", err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", db.Token))

	response, err := db.httpClient.Do(request)
	if err != nil {
		return []int64{}, fmt.Errorf("could not execute request: %w", err)
	}

	defer response.Body.Close()
	bodyContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []int64{}, fmt.Errorf("could not read response body: %w", err)
	}

	timestamps := make([]int64, 0)

	err = json.Unmarshal(bodyContent, &timestamps)
	if err != nil {
		return []int64{}, fmt.Errorf("could not unmarshal response: %w", err)
	}

	return timestamps, nil
}

func (db signalfxDB) getSegmentSpans(ctx context.Context, traceID string, timestamp int64) ([]signalFXSpan, error) {
	url := fmt.Sprintf("%s/v2/apm/trace/%s/%d", db.getURL(), traceID, timestamp)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []signalFXSpan{}, fmt.Errorf("could not create request: %w", err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", db.Token))

	response, err := db.httpClient.Do(request)
	if err != nil {
		return []signalFXSpan{}, fmt.Errorf("could not execute request: %w", err)
	}

	if response.StatusCode != 200 {
		return []signalFXSpan{}, nil
	}

	defer response.Body.Close()
	bodyContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []signalFXSpan{}, fmt.Errorf("could not read response body: %w", err)
	}

	spans := make([]signalFXSpan, 0)

	err = json.Unmarshal(bodyContent, &spans)
	if err != nil {
		return []signalFXSpan{}, fmt.Errorf("could not unmarshal response: %w", err)
	}

	return spans, nil
}

func convertSignalFXSpan(in signalFXSpan) traces.Span {
	attributes := make(traces.Attributes, 0)
	for name, value := range in.Tags {
		attributes[name] = value
	}

	for name, value := range in.ProcessTags {
		attributes[name] = value
	}

	attributes["parent_id"] = in.ParentID
	attributes["kind"] = attributes["span.kind"]
	delete(attributes, "span.kind")

	spanID, _ := trace.SpanIDFromHex(in.SpanID)
	startTime, _ := time.Parse(time.RFC3339, in.StartTime)
	endTime := startTime.Add(time.Duration(in.Duration) * time.Microsecond)

	return traces.Span{
		ID:         spanID,
		Name:       in.Name,
		StartTime:  startTime,
		EndTime:    endTime,
		Attributes: attributes,
	}
}

func newSignalFXDB(cfg config.SignalFXDataStoreConfig) (TraceDB, error) {
	return signalfxDB{
		Realm:      cfg.Realm,
		Token:      cfg.Token,
		httpClient: http.DefaultClient,
	}, nil
}

type signalFXSpan struct {
	TraceID     string            `json:"traceId"`
	SpanID      string            `json:"spanId"`
	ParentID    string            `json:"parentId"`
	Name        string            `json:"operationName"`
	StartTime   string            `json:"startTime"`
	Duration    int               `json:"durationMicros"`
	Tags        map[string]string `json:"tags"`
	ProcessTags map[string]string `json:"processTags"`
}
