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
	"log"
	"strings"
)

type elasticsearchDB struct {
	realTraceDB
	config *config.OpensearchDataStoreConfig
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

	_, err := db.GetTraceByID(ctx, trace.TraceID{}.String())
	if strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
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
		Index: []string{db.config.Index},
		Body:  content,
	}

	response, err := searchRequest.Do(ctx, db.client)
	if err != nil {
		return model.Trace{}, fmt.Errorf("could not execute search request: %w", err)
	}
	defer response.Body.Close()

	/*var r map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		response.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)*/

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return model.Trace{}, fmt.Errorf("could not read response body")
	}

	// log.Print(string(responseBody))

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

func newElasticSearchDB(cfg *config.OpensearchDataStoreConfig) (TraceDB, error) {
	cert := []byte("-----BEGIN CERTIFICATE-----\nMIIDSTCCAjGgAwIBAgIUKQkpbNBMH8ksjT5bX3Nc1bacRN4wDQYJKoZIhvcNAQEL\nBQAwNDEyMDAGA1UEAxMpRWxhc3RpYyBDZXJ0aWZpY2F0ZSBUb29sIEF1dG9nZW5l\ncmF0ZWQgQ0EwHhcNMjMwMTA5MTYzNjIwWhcNMjYwMTA4MTYzNjIwWjA0MTIwMAYD\nVQQDEylFbGFzdGljIENlcnRpZmljYXRlIFRvb2wgQXV0b2dlbmVyYXRlZCBDQTCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALh+5UdAvvNLBCETcSRRLUxV\nDtzkGEwZC1D5srK0yhyZLAO3+3uVYTazmrRAyBPtJ4IJ+WEmThp0DRuAcml3VtRC\nUvAmH0nRcpPd1D+X129xrn/pBZ2BNp3KDkY2+nKPXLHMaiTna+30jqSjShskLz6U\nrGzxJAJ9SCOAUPFajGearD3jSVYUxTshQ/k9Vvdb/2+968DqYBAFQIQlKYoajrn6\nun7ukUbGDcG3fvY+QijA3SGbAf+UIrTGV+BrYNVl6+ox9oVzNe0NM7iApqpFV1xC\niFf2APA0w0NM7uRR+rC6lhdhyBA7sY0PbN5GQdiLu31ZBMa9qXu8N22WMa9tiucC\nAwEAAaNTMFEwHQYDVR0OBBYEFFyiZYCC/d2ea7AnIk7F2NjUwVYKMB8GA1UdIwQY\nMBaAFFyiZYCC/d2ea7AnIk7F2NjUwVYKMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZI\nhvcNAQELBQADggEBAH3yYfb3SHzKig8KsWFEPifNQbqIHQHsEZOkScgdEErAuR/T\nrOkdKVju/O2ieh7ruwXYNOMttH9q735lFURJ0QuxNFqXf2L/05mWO5t3D9QCPv+2\nT1IhSX2+8sNAea5XAUPCWq7jq8IHbqIVgw6hABhXZe9hAjHgJIF1C2WNXiiYxrKr\nAYGjfPz575VDl6RyV/iISuBKr2trNM6V/mDqF4VwMyN00nNWWZMfhmhYR3HBBWhr\nhO7HLs5Tii5TnMjwGdl/zXmz7ASOpx0Nu0CrjxgBANhZQS1PkPd7t5zAntyjjvns\ndtnhPPVibWotnW47bJocEcW4Y62/TdBDW1ozhkY=\n-----END CERTIFICATE-----")

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		CACert:    cert,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create elasticsearch client: %w", err)
	}

	// Test: Get cluster info
	//
	var r map[string]interface{}
	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	return &elasticsearchDB{
		config: cfg,
		client: client,
	}, nil
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
	spanId, _ := trace.SpanIDFromHex("1234567891234567")
	// spanId, _ := trace.SpanIDFromHex(input["spanId"].(string))

	// startTime, _ := time.Parse(time.RFC3339, input["startTime"].(string))
	// endTime, _ := time.Parse(time.RFC3339, input["endTime"].(string))

	attributes := make(model.Attributes, 0)

	/*for attrName, attrValue := range input {
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
	attributes["parent_id"] = input["parentSpanId"].(string)*/

	return model.Span{
		ID: spanId,
		// Name: input["name"].(string),
		Name: "test",
		// StartTime:  startTime,
		// EndTime:    endTime,
		Attributes: attributes,
		Parent:     nil,
		Children:   []*model.Span{},
	}
}
