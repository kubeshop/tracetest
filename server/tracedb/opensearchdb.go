package tracedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
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

	if len(searchResponse.Hits.Hits) == 0 {
		return traces.Trace{}, ErrTraceNotFound
	}

	return traces.Trace{}, nil
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

type searchResponse struct {
	Hits searchHits `json:"hits"`
}

type searchHits struct {
	Hits []searchResult `json:"hits"`
}

type searchResult struct {
	Source map[string]interface{} `json:"_source"`
}
