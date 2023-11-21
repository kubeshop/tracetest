package tracedb

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

func elasticSearchDefaultPorts() []string {
	return []string{"9200"}
}

type elasticsearchDB struct {
	realTraceDB
	config *datastore.ElasticSearchConfig
	client *elasticsearch.Client
}

func (db *elasticsearchDB) Connect(ctx context.Context) error {
	return nil
}

func (db *elasticsearchDB) Close() error {
	// No need to close this db
	return nil
}

func (db *elasticsearchDB) GetEndpoints() string {
	return strings.Join(db.config.Addresses, ", ")
}

func (db *elasticsearchDB) TestConnection(ctx context.Context) model.ConnectionResult {
	tester := connection.NewTester(
		connection.WithPortLintingTest(connection.PortLinter("ElasticSearch", elasticSearchDefaultPorts(), db.config.Addresses...)),
		connection.WithConnectivityTest(connection.ConnectivityStep(model.ProtocolHTTP, db.config.Addresses...)),
		connection.WithPollingTest(connection.TracePollingTestStep(db)),
		connection.WithAuthenticationTest(connection.NewTestStep(func(ctx context.Context) (string, error) {
			_, err := getClusterInfo(db.client)
			if err != nil {
				return "Tracetest tried to execute an ElasticSearch API request but it failed due to authentication issues", err
			}

			return "Tracetest managed to authenticate with ElasticSearch", nil
		})),
	)

	return tester.TestConnection(ctx)
}

func (db *elasticsearchDB) Ready() bool {
	return db.client != nil
}

func (db *elasticsearchDB) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	if !db.Ready() {
		return traces.Trace{}, fmt.Errorf("ElasticSearch dataStore not ready")
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
		return traces.Trace{}, fmt.Errorf("could not execute search request: %w", err)
	}
	defer response.Body.Close()

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
		return traces.Trace{}, connection.ErrTraceNotFound
	}

	return convertElasticSearchFormatIntoTrace(traceID, searchResponse), nil
}

func newElasticSearchDB(cfg *datastore.ElasticSearchConfig) (TraceDB, error) {
	var caCert []byte
	if cfg.Certificate != "" {
		caCert = []byte(cfg.Certificate)
	}

	client, err := elasticsearch.NewClient(elasticsearch.Config{
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

func convertElasticSearchFormatIntoTrace(traceID string, searchResponse searchResponse) traces.Trace {
	spans := make([]traces.Span, 0)
	for _, result := range searchResponse.Hits.Results {
		span := convertElasticSearchSpanIntoSpan(result.Source)
		spans = append(spans, span)
	}

	return traces.NewTrace(traceID, spans)
}

func convertElasticSearchSpanIntoSpan(input map[string]interface{}) traces.Span {
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
	attributes := traces.NewAttributes()

	for attrName, attrValue := range flatInput {
		name := attrName
		name = strings.ReplaceAll(name, "transaction.", "")
		name = strings.ReplaceAll(name, "span.", "")
		attributes.Set(name, fmt.Sprintf("%v", attrValue))
	}

	// ParentId
	parentId := flatInput["parent.id"]
	if parentId != nil {
		attributes.Set(traces.TracetestMetadataFieldParentID, flatInput["parent.id"].(string))
	}

	return traces.Span{
		ID:         id,
		Name:       name,
		StartTime:  startTime,
		EndTime:    endTime,
		Attributes: attributes,
		Parent:     nil,
		Children:   []*traces.Span{},
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
