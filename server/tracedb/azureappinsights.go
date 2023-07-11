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
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
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

func NewAzureAppInsightsDB(config *datastore.AzureAppInsightsConfig) (TraceDB, error) {
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

type spanTable struct {
	rows []spanRow
}

func (st *spanTable) Spans() []spanRow {
	output := make([]spanRow, 0)
	for _, row := range st.rows {
		if row.Type() != "trace" {
			output = append(output, row)
		}
	}

	return output
}

func (st *spanTable) Events() []spanRow {
	output := make([]spanRow, 0)
	for _, row := range st.rows {
		if row.Type() == "trace" {
			output = append(output, row)
		}
	}

	return output
}

type spanRow struct {
	values map[string]any
}

func (sr *spanRow) Get(name string) any {
	return sr.values[name]
}

func (sr *spanRow) Type() string {
	return sr.values["itemType"].(string)
}

func (sr *spanRow) ParentID() string {
	return sr.values["operation_ParentId"].(string)
}

func (sr *spanRow) SpanID() string {
	return sr.values["id"].(string)
}

func newSpanTable(table *azquery.Table) spanTable {
	spanRows := make([]spanRow, 0, len(table.Rows))
	for _, row := range table.Rows {
		spanRows = append(spanRows, newSpanRow(row, table.Columns))
	}

	return spanTable{spanRows}
}

func newSpanRow(row azquery.Row, columns []*azquery.Column) spanRow {
	values := make(map[string]any)
	for i, column := range columns {
		name := *column.Name
		if value := row[i]; value != nil {
			values[name] = value
		}
	}

	return spanRow{values}
}

func parseAzureAppInsightsTrace(traceID string, table *azquery.Table) (model.Trace, error) {
	spans, err := parseSpans(table)
	if err != nil {
		return model.Trace{}, err
	}

	return model.NewTrace(traceID, spans), nil
}

func parseSpans(table *azquery.Table) ([]model.Span, error) {
	spanTable := newSpanTable(table)
	spanRows := spanTable.Spans()
	eventRows := spanTable.Events()

	spanEventsMap := make(map[string][]spanRow)
	for _, eventRow := range eventRows {
		spanEventsMap[eventRow.ParentID()] = append(spanEventsMap[eventRow.ParentID()], eventRow)
	}

	spanMap := make(map[string]*model.Span)
	for _, spanRow := range spanRows {
		span, err := parseRowToSpan(spanRow)
		if err != nil {
			return []model.Span{}, err
		}

		spanMap[span.ID.String()] = &span
	}

	for _, eventRow := range eventRows {
		parentSpan := spanMap[eventRow.ParentID()]
		event, err := parseEvent(eventRow)
		if err != nil {
			return []model.Span{}, err
		}

		parentSpan.Events = append(parentSpan.Events, event)
	}

	spans := make([]model.Span, 0, len(spanMap))
	for _, span := range spanMap {
		spans = append(spans, *span)
	}

	return spans, nil
}

func parseEvent(row spanRow) (model.SpanEvent, error) {
	event := model.SpanEvent{
		Name: row.Get("message").(string),
	}

	timestamp, err := time.Parse(time.RFC3339Nano, row.Get("timestamp").(string))
	if err != nil {
		return event, fmt.Errorf("could not parse event timestamp: %w", err)
	}

	event.Timestamp = timestamp

	attributes := make(model.Attributes, 0)
	rawAttributes := row.Get("customDimensions").(string)
	err = json.Unmarshal([]byte(rawAttributes), &attributes)
	if err != nil {
		return event, fmt.Errorf("could not unmarshal event attributes: %w", err)
	}

	event.Attributes = attributes

	return event, nil
}

func parseRowToSpan(row spanRow) (model.Span, error) {
	attributes := make(model.Attributes, 0)
	span := model.Span{
		Attributes: attributes,
	}
	var duration time.Duration

	for name, value := range row.values {
		switch name {
		case "id":
			err := parseSpanID(&span, value)
			if err != nil {
				return span, err
			}
		case "customDimensions":
			err := parseAttributes(&span, value)
			if err != nil {
				return span, err
			}
		case "operation_ParentId":
			err := parseParentID(&span, value)
			if err != nil {
				return span, err
			}
		case "name":
			err := parseName(&span, value)
			if err != nil {
				return span, err
			}
		case "timestamp":
			err := parseStartTime(&span, value)
			if err != nil {
				return span, err
			}
		case "duration":
			timeDuration, err := parseDuration(value)
			if err != nil {
				return span, err
			}

			duration = timeDuration
		}
	}

	span.EndTime = span.StartTime.Add(duration)
	return span, nil
}

func parseSpanID(span *model.Span, value any) error {
	spanID, err := trace.SpanIDFromHex(value.(string))
	if err != nil {
		return fmt.Errorf("failed to parse spanId: %w", err)
	}

	span.ID = spanID
	return nil
}

func parseAttributes(span *model.Span, value any) error {
	attributes := make(model.Attributes, 0)
	rawAttributes := value.(string)
	err := json.Unmarshal([]byte(rawAttributes), &attributes)
	if err != nil {
		return fmt.Errorf("failed to parse attributes: %w", err)
	}

	for key, value := range attributes {
		span.Attributes[key] = value
	}
	return nil
}

func parseParentID(span *model.Span, value any) error {
	rawParentID, ok := value.(string)
	if ok {
		span.Attributes[model.TracetestMetadataFieldParentID] = rawParentID
	} else {
		span.Attributes[model.TracetestMetadataFieldParentID] = ""
	}
	return nil
}

func parseName(span *model.Span, value any) error {
	rawName, ok := value.(string)
	if ok {
		span.Name = rawName
	} else {
		span.Name = ""
	}
	return nil
}

func parseStartTime(span *model.Span, value any) error {
	rawStartTime := value.(string)
	startTime, err := time.Parse(time.RFC3339Nano, rawStartTime)
	if err != nil {
		return fmt.Errorf("failed to parse startTime: %w", err)
	}

	span.StartTime = startTime
	return nil
}

func parseDuration(value any) (time.Duration, error) {
	rawDuration, ok := value.(float64)
	if !ok {
		return time.Duration(0), fmt.Errorf("failed to parse duration")
	}
	return time.Duration(rawDuration), nil
}

type tokenCredentials struct {
	accessToken string
}

func (c *tokenCredentials) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: c.accessToken}, nil
}
