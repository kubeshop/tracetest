package tracedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.opentelemetry.io/otel/trace"
)

type opensearchDb struct {
	config config.OpensearchDataStoreConfig
	client *opensearch.Client
}

func (db opensearchDb) Close() error {
	// No need to close this db
	return nil
}

func (db opensearchDb) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	content := strings.NewReader(fmt.Sprintf(`{
		"query": { "match": { "traceId": "%s" } }
	}`, traceID))

	searchRequest := opensearchapi.SearchRequest{
		Index: []string{db.config.Index},
		Body:  content,
	}

	response, err := searchRequest.Do(ctx, db.client)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not execute search request: %w", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not read response body")
	}

	var searchResponse searchResponse
	err = json.Unmarshal(responseBody, &searchResponse)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not unmarshal search response into struct: %w", err)
	}

	if len(searchResponse.Hits.Results) == 0 {
		return traces.Trace{}, ErrTraceNotFound
	}

	return convertOpensearchFormatIntoTrace(traceID, searchResponse), nil
}

func newOpenSearchDB(cfg config.OpensearchDataStoreConfig) (TraceDB, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create opensearch client: %w", err)
	}

	return opensearchDb{
		config: cfg,
		client: client,
	}, nil
}

func convertOpensearchFormatIntoTrace(traceID string, searchResponse searchResponse) traces.Trace {
	spans := make([]traces.Span, 0)
	for _, result := range searchResponse.Hits.Results {
		span := convertOpensearchSpanIntoSpan(result.Source)
		spans = append(spans, span)
	}

	return traces.New(traceID, spans)
}

func convertOpensearchSpanIntoSpan(input map[string]interface{}) traces.Span {
	spanId, _ := trace.SpanIDFromHex(input["spanId"].(string))

	startTime, _ := time.Parse(time.RFC3339, input["startTime"].(string))
	endTime, _ := time.Parse(time.RFC3339, input["endTime"].(string))

	attributes := make(traces.Attributes, 0)

	for attrName, attrValue := range input {
		if !strings.HasPrefix(attrName, "span.attributes.") && !strings.HasPrefix(attrName, "resource.attributes.") {
			// Not an attribute we care about
			continue
		}

		name := attrName
		name = strings.ReplaceAll(name, "span.attributes.", "")
		name = strings.ReplaceAll(name, "resource.attributes.", "")
		// Opensearch's data-prepper replaces "_" with "@". We have to revert it. Example:
		// "service_name" becomes "service@name"
		name = strings.ReplaceAll(name, "@", "_")
		attributes[name] = fmt.Sprintf("%v", attrValue)
	}

	attributes["kind"] = input["kind"].(string)
	attributes["parent_id"] = input["parentSpanId"].(string)

	return traces.Span{
		ID:         spanId,
		Name:       input["name"].(string),
		StartTime:  startTime,
		EndTime:    endTime,
		Attributes: attributes,
		Parent:     nil,
		Children:   []*traces.Span{},
	}
}

type searchResponse struct {
	Hits searchHits `json:"hits"`
}

type searchHits struct {
	Results []searchResult `json:"hits"`
}

type searchResult struct {
	Source map[string]interface{} `json:"_source"`
}
