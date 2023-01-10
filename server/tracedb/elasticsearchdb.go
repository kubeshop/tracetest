package tracedb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type elasticsearchDB struct {
	realTraceDB
	config *config.ElasticSearchDataStoreConfig
	client *elasticsearch.Client
}

func (db elasticsearchDB) Connect(ctx context.Context) error {
	return nil
}

func (db elasticsearchDB) Close() error {
	// No need to close this db
	return nil
}

func (db elasticsearchDB) TestConnection(ctx context.Context) ConnectionTestResult {
	addressesString := strings.Join(db.config.Addresses, ",")
	connectionTestResult := ConnectionTestResult{
		ConnectivityTestResult: ConnectionTestStepResult{
			OperationDescription: fmt.Sprintf(`Tracetest connected to "%s"`, addressesString),
		},
		AuthenticationTestResult: ConnectionTestStepResult{
			OperationDescription: `Tracetest managed to authenticate with ElasticSearch`,
		},
		TraceRetrivalTestResult: ConnectionTestStepResult{
			OperationDescription: `Tracetest was able to search for a trace using the ElasticSearch API`,
		},
	}

	for _, address := range db.config.Addresses {
		reachable, err := isReachable(address)

		if !reachable {
			return ConnectionTestResult{
				ConnectivityTestResult: ConnectionTestStepResult{
					OperationDescription: fmt.Sprintf(`Tracetest tried to connect to "%s" and failed`, address),
					Error:                err,
				},
			}
		}
	}

	_, err := getClusterInfo(db.client)
	if err != nil {
		return ConnectionTestResult{
			ConnectivityTestResult: connectionTestResult.ConnectivityTestResult,
			AuthenticationTestResult: ConnectionTestStepResult{
				OperationDescription: `Tracetest tried to execute an ElasticSearch API request but it failed due to authentication issues`,
				Error:                err,
			},
		}
	}

	_, err = db.GetTraceByID(ctx, trace.TraceID{}.String())
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
		return ConnectionTestResult{
			ConnectivityTestResult: connectionTestResult.ConnectivityTestResult,
			AuthenticationTestResult: ConnectionTestStepResult{
				OperationDescription: `Tracetest tried to execute an ElasticSearch API request but it failed due to authentication issues`,
				Error:                err,
			},
		}
	}

	if !errors.Is(err, ErrTraceNotFound) {
		return ConnectionTestResult{
			ConnectivityTestResult:   connectionTestResult.ConnectivityTestResult,
			AuthenticationTestResult: connectionTestResult.AuthenticationTestResult,
			TraceRetrivalTestResult: ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to fetch a trace from the ElasticSearch endpoint "%s" and got an error`, addressesString),
				Error:                err,
			},
		}
	}

	return connectionTestResult
}

func (db elasticsearchDB) Ready() bool {
	return db.client != nil
}

func (db elasticsearchDB) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	if !db.Ready() {
		return model.Trace{}, fmt.Errorf("ElasticSearch dataStore not ready")
	}
	content := strings.NewReader(fmt.Sprintf(`{
		"query": { "match": { "trace.id": "%s" } }
	}`, traceID))

	searchRequest := esapi.SearchRequest{
		Index:  []string{db.config.Index},
		Body:   content,
		Pretty: true,
	}

	response, err := searchRequest.Do(ctx, db.client)
	if err != nil {
		return model.Trace{}, fmt.Errorf("could not execute search request: %w", err)
	}
	defer response.Body.Close()

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
		return model.Trace{}, ErrTraceNotFound
	}

	return convertElasticSearchFormatIntoTrace(traceID, searchResponse), nil
}

func newElasticSearchDB(cfg *config.ElasticSearchDataStoreConfig) (TraceDB, error) {
	// cert := []byte("-----BEGIN CERTIFICATE-----\nMIIDSTCCAjGgAwIBAgIUKQkpbNBMH8ksjT5bX3Nc1bacRN4wDQYJKoZIhvcNAQEL\nBQAwNDEyMDAGA1UEAxMpRWxhc3RpYyBDZXJ0aWZpY2F0ZSBUb29sIEF1dG9nZW5l\ncmF0ZWQgQ0EwHhcNMjMwMTA5MTYzNjIwWhcNMjYwMTA4MTYzNjIwWjA0MTIwMAYD\nVQQDEylFbGFzdGljIENlcnRpZmljYXRlIFRvb2wgQXV0b2dlbmVyYXRlZCBDQTCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALh+5UdAvvNLBCETcSRRLUxV\nDtzkGEwZC1D5srK0yhyZLAO3+3uVYTazmrRAyBPtJ4IJ+WEmThp0DRuAcml3VtRC\nUvAmH0nRcpPd1D+X129xrn/pBZ2BNp3KDkY2+nKPXLHMaiTna+30jqSjShskLz6U\nrGzxJAJ9SCOAUPFajGearD3jSVYUxTshQ/k9Vvdb/2+968DqYBAFQIQlKYoajrn6\nun7ukUbGDcG3fvY+QijA3SGbAf+UIrTGV+BrYNVl6+ox9oVzNe0NM7iApqpFV1xC\niFf2APA0w0NM7uRR+rC6lhdhyBA7sY0PbN5GQdiLu31ZBMa9qXu8N22WMa9tiucC\nAwEAAaNTMFEwHQYDVR0OBBYEFFyiZYCC/d2ea7AnIk7F2NjUwVYKMB8GA1UdIwQY\nMBaAFFyiZYCC/d2ea7AnIk7F2NjUwVYKMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZI\nhvcNAQELBQADggEBAH3yYfb3SHzKig8KsWFEPifNQbqIHQHsEZOkScgdEErAuR/T\nrOkdKVju/O2ieh7ruwXYNOMttH9q735lFURJ0QuxNFqXf2L/05mWO5t3D9QCPv+2\nT1IhSX2+8sNAea5XAUPCWq7jq8IHbqIVgw6hABhXZe9hAjHgJIF1C2WNXiiYxrKr\nAYGjfPz575VDl6RyV/iISuBKr2trNM6V/mDqF4VwMyN00nNWWZMfhmhYR3HBBWhr\nhO7HLs5Tii5TnMjwGdl/zXmz7ASOpx0Nu0CrjxgBANhZQS1PkPd7t5zAntyjjvns\ndtnhPPVibWotnW47bJocEcW4Y62/TdBDW1ozhkY=\n-----END CERTIFICATE-----")
	var caCert []byte
	if cfg.Certificate != "" {
		caCert = []byte(cfg.Certificate)
	}

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		CACert:    caCert,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create elasticsearch client: %w", err)
	}

	return &elasticsearchDB{
		config: cfg,
		client: client,
	}, nil
}

func getClusterInfo(client *elasticsearch.Client) (string, error) {
	var r map[string]interface{}
	res, err := client.Info()
	if err != nil {
		return "", fmt.Errorf("error getting cluster info response: %s", err)
	}
	defer res.Body.Close()

	// Check response status
	if res.IsError() {
		return "", fmt.Errorf("error getting cluster info response status: %s", res.String())
	}
	// Deserialize the response into a map
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("error parsing cluster info response: %s", err)
	}

	// Return client version number
	info := fmt.Sprintf("Server: %s", r["version"].(map[string]interface{})["number"])
	return info, nil
}

func convertElasticSearchFormatIntoTrace(traceID string, searchResponse searchResponse) model.Trace {
	spans := make([]model.Span, 0)
	for _, result := range searchResponse.Hits.Results {
		span := convertElasticSearchSpanIntoSpan(result.Source)
		spans = append(spans, span)
	}

	return model.NewTrace(traceID, spans)
}

func convertElasticSearchSpanIntoSpan(input map[string]interface{}) model.Span {
	opts := &FlattenOptions{Delimiter: "."}
	flatInput, _ := flatten(opts.Prefix, 0, input, opts)

	// SpanId
	transactionId := flatInput["transaction.id"]
	spanId := flatInput["span.id"]
	var id trace.SpanID
	if transactionId != nil {
		id, _ = trace.SpanIDFromHex((transactionId).(string))
	}
	if spanId != nil {
		id, _ = trace.SpanIDFromHex((spanId).(string))
	}

	// SpanName
	transactionName := flatInput["transaction.name"]
	spanName := flatInput["span.name"]
	var name string
	if transactionName != nil {
		name = transactionName.(string)
	}
	if spanName != nil {
		name = spanName.(string)
	}

	// Duration
	transactionDuration := flatInput["transaction.duration.us"]
	spanDuration := flatInput["span.duration.us"]
	var duration float64
	if transactionDuration != nil {
		duration = transactionDuration.(float64)
	}
	if spanDuration != nil {
		duration = spanDuration.(float64)
	}

	// Timestamps
	startTime, _ := time.Parse(time.RFC3339, flatInput["@timestamp"].(string))
	endTime := startTime.Add(time.Microsecond * time.Duration(duration))

	// Attributes
	attributes := make(model.Attributes, 0)

	for attrName, attrValue := range flatInput {
		name := attrName
		name = strings.ReplaceAll(name, "transaction.", "")
		name = strings.ReplaceAll(name, "span.", "")
		attributes[name] = fmt.Sprintf("%v", attrValue)
	}

	// ParentId
	parentId := flatInput["parent.id"]
	if parentId != nil {
		attributes["parent_id"] = flatInput["parent.id"].(string)
	}

	return model.Span{
		ID:         id,
		Name:       name,
		StartTime:  startTime,
		EndTime:    endTime,
		Attributes: attributes,
		Parent:     nil,
		Children:   []*model.Span{},
	}
}

type FlattenOptions struct {
	Prefix    string
	Delimiter string
	Safe      bool
	MaxDepth  int
}

func flatten(prefix string, depth int, nested interface{}, opts *FlattenOptions) (flatmap map[string]interface{}, err error) {
	flatmap = make(map[string]interface{})

	switch nested := nested.(type) {
	case map[string]interface{}:
		if opts.MaxDepth != 0 && depth >= opts.MaxDepth {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, map[string]interface{}{}) {
			flatmap[prefix] = nested
			return
		}
		for k, v := range nested {
			// create new key
			newKey := k
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := flatten(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	case []interface{}:
		if opts.Safe {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, []interface{}{}) {
			flatmap[prefix] = nested
			return
		}
		for i, v := range nested {
			newKey := strconv.Itoa(i)
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := flatten(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	default:
		flatmap[prefix] = nested
	}
	return
}

func update(to map[string]interface{}, from map[string]interface{}) {
	for kt, vt := range from {
		to[kt] = vt
	}
}
