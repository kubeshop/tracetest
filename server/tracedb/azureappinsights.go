package tracedb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
	"go.opentelemetry.io/otel/trace"
)

var (
	azureLogUrl = "https://api.loganalytics.azure.com"
)

type azureAppInsightsDB struct {
	realTraceDB

	resourceArmId string
	credentials   azcore.TokenCredential
	client        *azquery.LogsClient
}

var _ TraceDB = &azureAppInsightsDB{}

func NewAzureAppInsightsDB(config *datastoreresource.AzureAppInsightsConfig) (TraceDB, error) {
	var credentials azcore.TokenCredential
	var err error
	if config.UseAzureActiveDirectoryAuth {
		credentials, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}
	} else {
		creds := []azcore.TokenCredential{
			&tokenCredentials{accessToken: config.AccessToken},
		}

		credentials, err = azidentity.NewChainedTokenCredential(creds, nil)
		if err != nil {
			return nil, err
		}
	}

	return &azureAppInsightsDB{
		resourceArmId: config.ResourceArmId,
		credentials:   credentials,
	}, nil
}

func (db *azureAppInsightsDB) Connect(ctx context.Context) error {
	client, err := azquery.NewLogsClient(db.credentials, nil)
	if err != nil {
		return err
	}

	db.client = client
	return nil
}

func (db *azureAppInsightsDB) Close() error {
	return nil
}

func (db *azureAppInsightsDB) Ready() bool {
	return db.credentials != nil && db.client != nil
}

func (db *azureAppInsightsDB) GetEndpoints() string {
	return azureLogUrl
}

func (db *azureAppInsightsDB) TestConnection(ctx context.Context) model.ConnectionResult {
	url := azureLogUrl
	tester := connection.NewTester(
		connection.WithConnectivityTest(connection.ConnectivityStep(model.ProtocolHTTP, url)),
		connection.WithPollingTest(connection.TracePollingTestStep(db)),
		connection.WithAuthenticationTest(connection.NewTestStep(func(ctx context.Context) (string, error) {
			_, err := db.GetTraceByID(ctx, db.GetTraceID().String())
			if err != nil && strings.Contains(strings.ToLower(err.Error()), "403") {
				return `Tracetest tried to execute an Azure API request but it failed due to authentication issues`, err
			}

			return "Tracetest managed to authenticate with the Azure Services", nil
		})),
	)

	return tester.TestConnection(ctx)
}

func (db *azureAppInsightsDB) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	query := fmt.Sprintf("union * | where operation_Id == '%s'", traceID)
	body := azquery.Body{
		Query: &query,
	}

	res, err := db.client.QueryResource(ctx, db.resourceArmId, body, nil)
	if err != nil {
		return model.Trace{}, err
	}

	table := res.Tables[0]
	if len(table.Rows) == 0 {
		return model.Trace{}, connection.ErrTraceNotFound
	}

	return parseAzureAppInsightsTrace(traceID, table)
}

type columnIndex map[string]int

func parseAzureAppInsightsTrace(traceID string, table *azquery.Table) (model.Trace, error) {
	columnIndex := mapColumnNames(table.Columns)
	spans := make([]model.Span, len(table.Rows))

	for i, row := range table.Rows {
		span, err := parseRowToSpan(row, columnIndex)
		if err != nil {
			return model.Trace{}, err
		}

		spans[i] = span
	}

	return model.NewTrace(traceID, spans), nil
}

var columnNamesMap = map[string]string{
	"traceId":    "operation_Id",
	"spanId":     "id",
	"parentId":   "operation_ParentId",
	"name":       "name",
	"attributes": "customDimensions",
	"startTime":  "timestamp",
	"duration":   "duration",
}

func parseRowToSpan(row azquery.Row, columnIndex columnIndex) (model.Span, error) {
	attributes := make(model.Attributes, 0)
	span := model.Span{
		Attributes: attributes,
	}
	var duration time.Duration

	for name, index := range columnIndex {
		switch name {
		case "spanId":
			err := parseSpanId(&span, row, index)
			if err != nil {
				return span, err
			}
		case "attributes":
			err := parseAttributes(&span, row, index)
			if err != nil {
				return span, err
			}
		case "parentId":
			err := parseParentId(&span, row, index)
			if err != nil {
				return span, err
			}
		case "name":
			err := parseName(&span, row, index)
			if err != nil {
				return span, err
			}
		case "startTime":
			err := parseStartTime(&span, row, index)
			if err != nil {
				return span, err
			}
		case "duration":
			timeDuration, err := parseDuration(row, index)
			if err != nil {
				return span, err
			}

			duration = timeDuration
		}
	}

	span.EndTime = span.StartTime.Add(duration)
	return span, nil
}

func parseSpanId(span *model.Span, row azquery.Row, index int) error {
	if index == -1 {
		return fmt.Errorf("spanId column not found")
	}

	rawSpanId := row[index].(string)
	spanId, err := trace.SpanIDFromHex(rawSpanId)
	if err != nil {
		return fmt.Errorf("failed to parse spanId: %w", err)
	}

	span.ID = spanId
	return nil
}

func parseAttributes(span *model.Span, row azquery.Row, index int) error {
	if index == -1 {
		return fmt.Errorf("attributes column not found")
	}

	attributes := make(model.Attributes, 0)
	rawAttributes := row[index].(string)
	err := json.Unmarshal([]byte(rawAttributes), &attributes)
	if err != nil {
		return fmt.Errorf("failed to parse attributes: %w", err)
	}

	for key, value := range attributes {
		span.Attributes[key] = value
	}
	return nil
}

func parseParentId(span *model.Span, row azquery.Row, index int) error {
	if index == -1 {
		return fmt.Errorf("parentId column not found")
	}

	rawParentId, ok := row[index].(string)
	if ok {
		span.Attributes[string(model.TracetestMetadataFieldParentID)] = rawParentId
	} else {
		span.Attributes[string(model.TracetestMetadataFieldParentID)] = ""
	}
	return nil
}

func parseName(span *model.Span, row azquery.Row, index int) error {
	if index == -1 {
		return fmt.Errorf("name column not found")
	}

	rawName, ok := row[index].(string)
	if ok {
		span.Name = rawName
	} else {
		span.Name = ""
	}
	return nil
}

func parseStartTime(span *model.Span, row azquery.Row, index int) error {
	if index == -1 {
		return fmt.Errorf("startTime column not found")
	}

	rawStartTime := row[index].(string)
	startTime, err := time.Parse(time.RFC3339Nano, rawStartTime)
	if err != nil {
		return fmt.Errorf("failed to parse startTime: %w", err)
	}

	span.StartTime = startTime
	return nil
}

func parseDuration(row azquery.Row, index int) (time.Duration, error) {
	if index == -1 {
		return time.Duration(0), fmt.Errorf("duration column not found")
	}

	rawDuration, ok := row[index].(float64)
	if !ok {
		return time.Duration(0), fmt.Errorf("failed to parse duration")
	}
	return time.Duration(rawDuration), nil
}

func mapColumnNames(columns []*azquery.Column) columnIndex {
	columnIndex := columnIndex{
		"traceId":    -1,
		"parentId":   -1,
		"name":       -1,
		"attributes": -1,
		"startTime":  -1,
		"duration":   -1,
	}

	for name, azureName := range columnNamesMap {
		columnIndex[name] = findColumnByName(columns, azureName)
	}

	return columnIndex
}

func findColumnByName(columns []*azquery.Column, name string) int {
	for i, column := range columns {
		if *column.Name == name {
			return i
		}
	}

	return -1
}

type tokenCredentials struct {
	accessToken string
}

func (c *tokenCredentials) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: c.accessToken}, nil
}
