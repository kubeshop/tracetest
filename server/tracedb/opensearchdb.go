package tracedb

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.opentelemetry.io/otel/trace"
)

type opensearchDB struct {
	realTraceDB
	config *model.ElasticSearchDataStoreConfig
	client *opensearch.Client
}

func (db opensearchDB) Connect(ctx context.Context) error {
	return nil
}

func (db opensearchDB) Close() error {
	// No need to close this db
	return nil
}

func (db opensearchDB) TestConnection(ctx context.Context) connection.ConnectionTestResult {
	addressesString := strings.Join(db.config.Addresses, ",")
	connectionTestResult := connection.ConnectionTestResult{
		ConnectivityTestResult: connection.ConnectionTestStepResult{
			OperationDescription: fmt.Sprintf(`Tracetest connected to "%s"`, addressesString),
		},
		AuthenticationTestResult: connection.ConnectionTestStepResult{
			OperationDescription: `Tracetest managed to authenticate with OpenSearch`,
		},
		TraceRetrievalTestResult: connection.ConnectionTestStepResult{
			OperationDescription: `Tracetest was able to search for a trace using the OpenSearch API`,
		},
	}

	for _, address := range db.config.Addresses {
		reachable, err := connection.IsReachable(address)

		if !reachable {
			return connection.ConnectionTestResult{
				ConnectivityTestResult: connection.ConnectionTestStepResult{
					OperationDescription: fmt.Sprintf(`Tracetest tried to connect to "%s" and failed`, address),
					Error:                err,
				},
			}
		}
	}

	_, err := db.GetTraceByID(ctx, trace.TraceID{}.String())
	if strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
		return connection.ConnectionTestResult{
			ConnectivityTestResult: connectionTestResult.ConnectivityTestResult,
			AuthenticationTestResult: connection.ConnectionTestStepResult{
				OperationDescription: `Tracetest tried to execute an OpenSearch API request but it failed due to authentication issues`,
				Error:                err,
			},
		}
	}

	if !errors.Is(err, connection.ErrTraceNotFound) {
		return connection.ConnectionTestResult{
			ConnectivityTestResult:   connectionTestResult.ConnectivityTestResult,
			AuthenticationTestResult: connectionTestResult.AuthenticationTestResult,
			TraceRetrievalTestResult: connection.ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to fetch a trace from the OpenSearch endpoint "%s" and got an error`, addressesString),
				Error:                err,
			},
		}
	}

	return connectionTestResult
}

func (db opensearchDB) Ready() bool {
	return db.client != nil
}

func (db opensearchDB) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	if !db.Ready() {
		return model.Trace{}, fmt.Errorf("OpenSearch dataStore not ready")
	}
	content := strings.NewReader(fmt.Sprintf(`{
		"query": { "match": { "traceId": "%s" } }
	}`, traceID))

	searchRequest := opensearchapi.SearchRequest{
		Index: []string{db.config.Index},
		Body:  content,
	}

	response, err := searchRequest.Do(ctx, db.client)
	if err != nil {
		return model.Trace{}, fmt.Errorf("could not execute search request: %w", err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return model.Trace{}, fmt.Errorf("could not read response body")
	}

	var searchResponse searchResponse
	err = json.Unmarshal(responseBody, &searchResponse)
	if err != nil {
		return model.Trace{}, fmt.Errorf("could not unmarshal search response into struct: %w", err)
	}

	if len(searchResponse.Hits.Results) == 0 {
		return model.Trace{}, connection.ErrTraceNotFound
	}

	return convertOpensearchFormatIntoTrace(traceID, searchResponse), nil
}

func newOpenSearchDB(cfg *model.ElasticSearchDataStoreConfig) (TraceDB, error) {
	var caCert []byte
	if cfg.Certificate != "" {
		caCert = []byte(cfg.Certificate)
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		CACert:    caCert,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.InsecureSkipVerify,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("could not create opensearch client: %w", err)
	}

	return &opensearchDB{
		config: cfg,
		client: client,
	}, nil
}

func convertOpensearchFormatIntoTrace(traceID string, searchResponse searchResponse) model.Trace {
	spans := make([]model.Span, 0)
	for _, result := range searchResponse.Hits.Results {
		span := convertOpensearchSpanIntoSpan(result.Source)
		spans = append(spans, span)
	}

	return model.NewTrace(traceID, spans)
}

func convertOpensearchSpanIntoSpan(input map[string]interface{}) model.Span {
	spanId, _ := trace.SpanIDFromHex(input["spanId"].(string))

	startTime, _ := time.Parse(time.RFC3339, input["startTime"].(string))
	endTime, _ := time.Parse(time.RFC3339, input["endTime"].(string))

	attributes := make(model.Attributes, 0)

	for attrName, attrValue := range input {
		if !strings.HasPrefix(attrName, "span.attributes.") && !strings.HasPrefix(attrName, "resource.attributes.") {
			// Not an attribute we care about
			continue
		}

		name := attrName
		name = strings.ReplaceAll(name, "span.attributes.", "")
		name = strings.ReplaceAll(name, "resource.attributes.", "")
		// Opensearch's data-prepper replaces "." with "@". We have to revert it. Example:
		// "service.name" becomes "service@name"
		name = strings.ReplaceAll(name, "@", ".")
		attributes[name] = fmt.Sprintf("%v", attrValue)
	}

	attributes["kind"] = input["kind"].(string)
	attributes["parent_id"] = input["parentSpanId"].(string)

	return model.Span{
		ID:         spanId,
		Name:       input["name"].(string),
		StartTime:  startTime,
		EndTime:    endTime,
		Attributes: attributes,
		Parent:     nil,
		Children:   []*model.Span{},
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
