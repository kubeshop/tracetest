package tracedb

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/tracedb/datastore"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	"go.opentelemetry.io/otel/trace"
)

type awsxrayDB struct {
	realTraceDB

	credentials    *credentials.Credentials
	session        *session.Session
	region         string
	service        *xray.XRay
	useDefaultAuth bool
}

func NewAwsXRayDB(cfg *datastore.AWSXRayConfig) (TraceDB, error) {
	sessionCredentials := credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken)

	return &awsxrayDB{
		credentials:    sessionCredentials,
		region:         cfg.Region,
		useDefaultAuth: cfg.UseDefaultAuth,
	}, nil
}

func (db *awsxrayDB) GetTraceID() trace.TraceID {
	var r [16]byte
	epoch := time.Now().Unix()
	binary.BigEndian.PutUint32(r[0:4], uint32(epoch))
	_, err := rand.Read(r[4:])
	if err != nil {
		panic(err)
	}

	return trace.TraceID(r)
}

func (db *awsxrayDB) Connect(ctx context.Context) error {
	awsConfig := &aws.Config{}

	if db.useDefaultAuth {
		awsConfig = aws.NewConfig().WithRegion(db.region)
	} else {
		awsConfig = &aws.Config{
			Region:      &db.region,
			Credentials: db.credentials,
		}
	}
	sess, err := session.NewSession(awsConfig)

	if err != nil {
		return err
	}

	db.service = xray.New(sess)
	db.session = sess

	return nil
}

func (db *awsxrayDB) Ready() bool {
	return db.service != nil
}

func (db *awsxrayDB) Close() error {
	// Doesn't need to be closed
	return nil
}

func (db *awsxrayDB) GetEndpoints() string {
	return fmt.Sprintf("xray.%s.amazonaws.com:443", db.region)
}

func (db *awsxrayDB) TestConnection(ctx context.Context) model.ConnectionResult {
	url := fmt.Sprintf("xray.%s.amazonaws.com:443", db.region)
	tester := connection.NewTester(
		connection.WithConnectivityTest(connection.ConnectivityStep(model.ProtocolHTTP, url)),
		connection.WithPollingTest(connection.TracePollingTestStep(db)),
		connection.WithAuthenticationTest(connection.NewTestStep(func(ctx context.Context) (string, error) {
			_, err := db.GetTraceByID(ctx, db.GetTraceID().String())
			if err != nil && strings.Contains(strings.ToLower(err.Error()), "403") {
				return `Tracetest tried to execute an AWS XRay API request but it failed due to authentication issues`, err
			}

			return "Tracetest managed to authenticate with the AWS Services", nil
		})),
	)

	return tester.TestConnection(ctx)
}

func (db *awsxrayDB) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	hexTraceID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return model.Trace{}, err
	}

	parsedTraceID, err := convertToAmazonTraceID(hexTraceID)
	if err != nil {
		return model.Trace{}, err
	}

	res, err := db.service.BatchGetTraces(&xray.BatchGetTracesInput{
		TraceIds: []*string{&parsedTraceID},
	})

	if err != nil {
		return model.Trace{}, err
	}

	if len(res.Traces) == 0 {
		return model.Trace{}, connection.ErrTraceNotFound
	}

	return parseXRayTrace(traceID, res.Traces[0])
}

func parseXRayTrace(traceID string, rawTrace *xray.Trace) (model.Trace, error) {
	if len(rawTrace.Segments) == 0 {
		return model.Trace{}, nil
	}

	spans := []model.Span{}

	for _, segment := range rawTrace.Segments {
		newSpans, err := parseSegmentToSpans([]byte(*segment.Document), traceID)

		if err != nil {
			return model.Trace{}, err
		}

		spans = append(spans, newSpans...)
	}

	return model.NewTrace(traceID, spans), nil
}

func parseSegmentToSpans(rawSeg []byte, traceID string) ([]model.Span, error) {
	var seg segment
	err := json.Unmarshal(rawSeg, &seg)
	if err != nil {
		return []model.Span{}, err
	}

	err = seg.Validate()
	if err != nil {
		return []model.Span{}, err
	}

	return segToSpans(seg, traceID, nil)
}

func segToSpans(seg segment, traceID string, parent *model.Span) ([]model.Span, error) {
	span, err := generateSpan(&seg, parent)
	if err != nil {
		return []model.Span{}, err
	}

	spans := []model.Span{span}

	for _, s := range seg.Subsegments {
		nestedSpans, err := segToSpans(s, traceID, &span)

		if err != nil {
			return spans, err
		}

		spans = append(spans, nestedSpans...)
	}

	return spans, nil
}

func generateSpan(seg *segment, parent *model.Span) (model.Span, error) {
	attributes := make(model.Attributes, 0)
	span := model.Span{
		Parent: parent,
		Name:   *seg.Name,
	}

	if seg.ParentID != nil {
		parentID, err := decodeXRaySpanID(seg.ParentID)
		if err != nil {
			return span, err
		}

		attributes[model.TracetestMetadataFieldParentID] = parentID.String()
	} else if parent != nil {
		attributes[model.TracetestMetadataFieldParentID] = parent.ID.String()
	}

	// decode span id
	spanID, err := decodeXRaySpanID(seg.ID)
	if err != nil {
		return span, err
	}
	span.ID = spanID

	err = addNamespace(seg, attributes)
	if err != nil {
		return model.Span{}, err
	}

	span.StartTime = floatSecToTime(seg.StartTime)
	if seg.EndTime != nil {
		span.EndTime = floatSecToTime(seg.EndTime)
	}

	if seg.InProgress != nil {
		attributes[AWSXRayInProgressAttribute] = strconv.FormatBool(*seg.InProgress)
	}

	attributes.SetPointerValue(conventions.AttributeEnduserID, seg.User)
	addHTTP(seg, attributes)
	addAWSToSpan(seg.AWS, attributes)
	err = addSQLToSpan(seg.SQL, attributes)
	if err != nil {
		return model.Span{}, err
	}

	if seg.Traced != nil {
		attributes[AWSXRayTracedAttribute] = strconv.FormatBool(*seg.Traced)
	}

	addAnnotations(seg.Annotations, attributes)
	addMetadata(seg.Metadata, attributes)

	// this generates an event that we don't support yet
	// addCause(seg, span)

	span.Attributes = attributes
	return span, nil
}

const (
	validAWSNamespace    = "aws"
	validRemoteNamespace = "remote"
)

func addNamespace(seg *segment, attributes model.Attributes) error {
	if seg.Namespace != nil {
		switch *seg.Namespace {
		case validAWSNamespace:
			attributes.SetPointerValue(AWSServiceAttribute, seg.Name)

		case validRemoteNamespace:
			// no op
		default:
			return fmt.Errorf("unexpected namespace: %s", *seg.Namespace)
		}
		return nil
	}

	return nil
}

func addHTTP(seg *segment, attributes model.Attributes) {
	if seg.HTTP == nil {
		return
	}

	if req := seg.HTTP.Request; req != nil {
		attributes.SetPointerValue(conventions.AttributeHTTPMethod, req.Method)
		attributes.SetPointerValue(conventions.AttributeHTTPClientIP, req.ClientIP)
		attributes.SetPointerValue(conventions.AttributeHTTPUserAgent, req.UserAgent)
		attributes.SetPointerValue(conventions.AttributeHTTPURL, req.URL)

		if req.XForwardedFor != nil {
			attributes[AWSXRayXForwardedForAttribute] = strconv.FormatBool(*req.XForwardedFor)
		}
	}

	if resp := seg.HTTP.Response; resp != nil {
		if resp.status != nil {
			attributes[conventions.AttributeHTTPStatusCode] = fmt.Sprintf("%v", *resp.status)
		}

		switch val := resp.contentLength.(type) {
		case string:
			attributes[conventions.AttributeHTTPResponseContentLength] = val
		case float64:
			lengthPointer := int64(val)
			attributes[conventions.AttributeHTTPResponseContentLength] = fmt.Sprintf("%v", lengthPointer)
		}
	}
}

func addAWSToSpan(aws *aWSData, attrs model.Attributes) {
	if aws != nil {
		attrs.SetPointerValue(AWSAccountAttribute, aws.AccountID)
		attrs.SetPointerValue(AWSOperationAttribute, aws.Operation)
		attrs.SetPointerValue(AWSRegionAttribute, aws.RemoteRegion)
		attrs.SetPointerValue(AWSRequestIDAttribute, aws.RequestID)
		attrs.SetPointerValue(AWSQueueURLAttribute, aws.QueueURL)
		attrs.SetPointerValue(AWSTableNameAttribute, aws.TableName)

		if aws.Retries != nil {
			attrs[AWSXrayRetriesAttribute] = fmt.Sprintf("%v", *aws.Retries)
		}
	}
}

func addSQLToSpan(sql *sQLData, attrs model.Attributes) error {
	if sql == nil {
		return nil
	}

	if sql.URL != nil {
		dbURL, dbName, err := splitSQLURL(*sql.ConnectionString)
		if err != nil {
			return err
		}

		attrs[conventions.AttributeDBConnectionString] = dbURL
		attrs[conventions.AttributeDBName] = dbName
	}
	// not handling sql.ConnectionString for now because the X-Ray exporter
	// does not support it
	attrs.SetPointerValue(conventions.AttributeDBSystem, sql.DatabaseType)
	attrs.SetPointerValue(conventions.AttributeDBStatement, sql.SanitizedQuery)
	attrs.SetPointerValue(conventions.AttributeDBUser, sql.User)
	return nil
}

func addAnnotations(annos map[string]interface{}, attrs model.Attributes) {
	if len(annos) > 0 {
		for k, v := range annos {
			switch t := v.(type) {
			case int:
				attrs[k] = fmt.Sprintf("%v", t)
			case int32:
				attrs[k] = fmt.Sprintf("%v", t)
			case int64:
				attrs[k] = fmt.Sprintf("%v", t)
			case string:
				attrs[k] = t
			case bool:
				attrs[k] = strconv.FormatBool(t)
			case float32:
				attrs[k] = fmt.Sprintf("%v", t)
			case float64:
				attrs[k] = fmt.Sprintf("%v", t)
			default:
			}
		}
	}
}

func addMetadata(meta map[string]map[string]interface{}, attrs model.Attributes) error {
	for k, v := range meta {
		val, err := json.Marshal(v)
		if err != nil {
			return err
		}
		attrs[AWSXraySegmentMetadataAttributePrefix+k] = string(val)
	}
	return nil
}

// SQL URL is of the format: protocol+transport://host:port/dbName?queryParam
var re = regexp.MustCompile(`^(.+\/\/.+)\/([^\?]+)\??.*$`)

const (
	dbURLI  = 1
	dbNameI = 2
)

func splitSQLURL(rawURL string) (string, string, error) {
	m := re.FindStringSubmatch(rawURL)
	if len(m) == 0 {
		return "", "", fmt.Errorf(
			"failed to parse out the database name in the \"sql.url\" field, rawUrl: %s",
			rawURL,
		)
	}
	return m[dbURLI], m[dbNameI], nil
}

func floatSecToTime(epochSec *float64) time.Time {
	timestamp := (*epochSec) * float64(time.Second)
	return time.Unix(0, int64(timestamp)).UTC()
}

const (
	traceIDLength    = 35 // fixed length of aws trace id
	identifierOffset = 11 // offset of identifier within traceID
)

func convertToAmazonTraceID(traceID trace.TraceID) (string, error) {
	const (
		maxAge  = 60 * 60 * 24 * 28
		maxSkew = 60 * 5
	)

	var (
		content      = [traceIDLength]byte{}
		epochNow     = time.Now().Unix()
		traceIDBytes = traceID
		epoch        = int64(binary.BigEndian.Uint32(traceIDBytes[0:4]))
		b            = [4]byte{}
	)

	delta := epochNow - epoch
	if delta > maxAge || delta < -maxSkew {
		return "", fmt.Errorf("invalid xray traceid: %s", traceID)
	}

	binary.BigEndian.PutUint32(b[0:4], uint32(epoch))

	content[0] = '1'
	content[1] = '-'
	hex.Encode(content[2:10], b[0:4])
	content[10] = '-'
	hex.Encode(content[identifierOffset:], traceIDBytes[4:16]) // overwrite with identifier

	return string(content[0:traceIDLength]), nil
}

func decodeXRaySpanID(spanID *string) (trace.SpanID, error) {
	sid := [8]byte{}
	if spanID == nil {
		return sid, errors.New("spanid is null")
	}
	if len(*spanID) != 16 {
		return sid, errors.New("spanID length is wrong")
	}
	_, err := hex.Decode(sid[:], []byte(*spanID))
	return sid, err
}

const (
	// TypeStr is the type and ingest format of this receiver
	TypeStr = "awsxray"
)

type CauseType int

const (
	CauseTypeExceptionID CauseType = iota + 1
	CauseTypeObject
)

type segment struct {
	// Required fields for both segment and subsegments
	Name      *string  `json:"name"`
	ID        *string  `json:"id"`
	StartTime *float64 `json:"start_time"`

	// Segment-only optional fields
	Service     *serviceData `json:"service,omitempty"`
	Origin      *string      `json:"origin,omitempty"`
	User        *string      `json:"user,omitempty"`
	ResourceARN *string      `json:"resource_arn,omitempty"`

	// Optional fields for both Segment and subsegments
	TraceID     *string                           `json:"trace_id,omitempty"`
	EndTime     *float64                          `json:"end_time,omitempty"`
	InProgress  *bool                             `json:"in_progress,omitempty"`
	HTTP        *hTTPData                         `json:"http,omitempty"`
	Fault       *bool                             `json:"fault,omitempty"`
	Error       *bool                             `json:"error,omitempty"`
	Throttle    *bool                             `json:"throttle,omitempty"`
	Cause       *causeData                        `json:"cause,omitempty"`
	AWS         *aWSData                          `json:"aws,omitempty"`
	Annotations map[string]interface{}            `json:"annotations,omitempty"`
	Metadata    map[string]map[string]interface{} `json:"metadata,omitempty"`
	Subsegments []segment                         `json:"subsegments,omitempty"`

	// (for both embedded and independent) subsegment-only (optional) fields.
	// Please refer to https://docs.aws.amazon.com/xray/latest/devguide/xray-api-segmentdocuments.html#api-segmentdocuments-subsegments
	// for more information on subsegment.
	Namespace    *string  `json:"namespace,omitempty"`
	ParentID     *string  `json:"parent_id,omitempty"`
	Type         *string  `json:"type,omitempty"`
	PrecursorIDs []string `json:"precursor_ids,omitempty"`
	Traced       *bool    `json:"traced,omitempty"`
	SQL          *sQLData `json:"sql,omitempty"`
}

// Validate checks whether the segment is valid or not
func (s *segment) Validate() error {
	if s.Name == nil {
		return errors.New(`segment "name" can not be nil`)
	}

	if s.ID == nil {
		return errors.New(`segment "id" can not be nil`)
	}

	if s.StartTime == nil {
		return errors.New(`segment "start_time" can not be nil`)
	}

	// it's ok for embedded subsegments to not have trace_id
	// but the root segment and independent subsegments must all
	// have trace_id.
	if s.TraceID == nil {
		return errors.New(`segment "trace_id" can not be nil`)
	}

	return nil
}

type aWSData struct {
	// Segment-only
	Beanstalk *beanstalkMetadata `json:"elastic_beanstalk,omitempty"`
	CWLogs    []logGroupMetadata `json:"cloudwatch_logs,omitempty"`
	ECS       *eCSMetadata       `json:"ecs,omitempty"`
	EC2       *eC2Metadata       `json:"ec2,omitempty"`
	EKS       *eKSMetadata       `json:"eks,omitempty"`
	XRay      *xRayMetaData      `json:"xray,omitempty"`

	// For both segment and subsegments
	AccountID    *string `json:"account_id,omitempty"`
	Operation    *string `json:"operation,omitempty"`
	RemoteRegion *string `json:"region,omitempty"`
	RequestID    *string `json:"request_id,omitempty"`
	QueueURL     *string `json:"queue_url,omitempty"`
	TableName    *string `json:"table_name,omitempty"`
	Retries      *int64  `json:"retries,omitempty"`
}

type eC2Metadata struct {
	InstanceID       *string `json:"instance_id"`
	AvailabilityZone *string `json:"availability_zone"`
	InstanceSize     *string `json:"instance_size"`
	AmiID            *string `json:"ami_id"`
}

type eCSMetadata struct {
	ContainerName    *string `json:"container,omitempty"`
	ContainerID      *string `json:"container_id,omitempty"`
	TaskArn          *string `json:"task_arn,omitempty"`
	TaskFamily       *string `json:"task_family,omitempty"`
	ClusterArn       *string `json:"cluster_arn,omitempty"`
	ContainerArn     *string `json:"container_arn,omitempty"`
	AvailabilityZone *string `json:"availability_zone,omitempty"`
	LaunchType       *string `json:"launch_type,omitempty"`
}

// BeanstalkMetadata represents the Elastic Beanstalk environment metadata field
type beanstalkMetadata struct {
	Environment  *string `json:"environment_name"`
	VersionLabel *string `json:"version_label"`
	DeploymentID *int64  `json:"deployment_id"`
}

// EKSMetadata represents the EKS metadata field
type eKSMetadata struct {
	ClusterName *string `json:"cluster_name"`
	Pod         *string `json:"pod"`
	ContainerID *string `json:"container_id"`
}

// LogGroupMetadata represents a single CloudWatch Log Group
type logGroupMetadata struct {
	LogGroup *string `json:"log_group"`
	Arn      *string `json:"arn,omitempty"`
}

// CauseData is the container that contains the `cause` field
type causeData struct {
	Type CauseType `json:"-"`
	// it will contain one of ExceptionID or (WorkingDirectory, Paths, Exceptions)
	ExceptionID *string `json:"-"`

	causeObject
}

type causeObject struct {
	WorkingDirectory *string     `json:"working_directory,omitempty"`
	Paths            []string    `json:"paths,omitempty"`
	Exceptions       []exception `json:"exceptions,omitempty"`
}

// UnmarshalJSON is the custom unmarshaller for the cause field
func (c *causeData) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &c.causeObject)
	if err == nil {
		c.Type = CauseTypeObject
		return nil
	}
	rawStr := string(data)
	if len(rawStr) > 0 && (rawStr[0] != '"' || rawStr[len(rawStr)-1] != '"') {
		return fmt.Errorf("the value assigned to the `cause` field does not appear to be a string: %v", data)
	}
	exceptionID := rawStr[1 : len(rawStr)-1]

	c.Type = CauseTypeExceptionID
	c.ExceptionID = &exceptionID
	return nil
}

// Exception represents an exception occurred
type exception struct {
	ID        *string      `json:"id,omitempty"`
	Message   *string      `json:"message,omitempty"`
	Type      *string      `json:"type,omitempty"`
	Remote    *bool        `json:"remote,omitempty"`
	Truncated *int64       `json:"truncated,omitempty"`
	Skipped   *int64       `json:"skipped,omitempty"`
	Cause     *string      `json:"cause,omitempty"`
	Stack     []stackFrame `json:"stack,omitempty"`
}

// StackFrame represents a frame in the stack when an exception occurred
type stackFrame struct {
	Path  *string `json:"path,omitempty"`
	Line  *int    `json:"line,omitempty"`
	Label *string `json:"label,omitempty"`
}

// HTTPData provides the shape for unmarshalling request and response fields.
type hTTPData struct {
	Request  *requestData  `json:"request,omitempty"`
	Response *responseData `json:"response,omitempty"`
}

// RequestData provides the shape for unmarshalling the request field.
type requestData struct {
	// Available in segment
	XForwardedFor *bool `json:"x_forwarded_for,omitempty"`

	// Available in both segment and subsegments
	Method    *string `json:"method,omitempty"`
	URL       *string `json:"url,omitempty"`
	UserAgent *string `json:"user_agent,omitempty"`
	ClientIP  *string `json:"client_ip,omitempty"`
}

// ResponseData provides the shape for unmarshalling the response field.
type responseData struct {
	status        *int64      `json:"status,omitempty"`
	contentLength interface{} `json:"content_length,omitempty"`
}

// ECSData provides the shape for unmarshalling the ecs field.
type eCSData struct {
	Container *string `json:"container"`
}

// EC2Data provides the shape for unmarshalling the ec2 field.
type eC2Data struct {
	InstanceID       *string `json:"instance_id"`
	AvailabilityZone *string `json:"availability_zone"`
}

// ElasticBeanstalkData provides the shape for unmarshalling the elastic_beanstalk field.
type elasticBeanstalkData struct {
	EnvironmentName *string `json:"environment_name"`
	VersionLabel    *string `json:"version_label"`
	DeploymentID    *int    `json:"deployment_id"`
}

// XRayMetaData provides the shape for unmarshalling the xray field
type xRayMetaData struct {
	SDK                 *string `json:"sdk,omitempty"`
	SDKVersion          *string `json:"sdk_version,omitempty"`
	AutoInstrumentation *bool   `json:"auto_instrumentation"`
}

// SQLData provides the shape for unmarshalling the sql field.
type sQLData struct {
	ConnectionString *string `json:"connection_string,omitempty"`
	URL              *string `json:"url,omitempty"` // protocol://host[:port]/database
	SanitizedQuery   *string `json:"sanitized_query,omitempty"`
	DatabaseType     *string `json:"database_type,omitempty"`
	DatabaseVersion  *string `json:"database_version,omitempty"`
	DriverVersion    *string `json:"driver_version,omitempty"`
	User             *string `json:"user,omitempty"`
	Preparation      *string `json:"preparation,omitempty"` // "statement" / "call"
}

// ServiceData provides the shape for unmarshalling the service field.
type serviceData struct {
	Version         *string `json:"version,omitempty"`
	CompilerVersion *string `json:"compiler_version,omitempty"`
	Compiler        *string `json:"compiler,omitempty"`
}

const (
	AWSOperationAttribute = "aws.operation"
	AWSAccountAttribute   = "aws.account_id"
	AWSRegionAttribute    = "aws.region"
	AWSRequestIDAttribute = "aws.request_id"
	// Currently different instrumentation uses different tag formats.
	// TODO(anuraaga): Find current instrumentation and consolidate.
	AWSRequestIDAttribute2 = "aws.requestId"
	AWSQueueURLAttribute   = "aws.queue_url"
	AWSQueueURLAttribute2  = "aws.queue.url"
	AWSServiceAttribute    = "aws.service"
	AWSTableNameAttribute  = "aws.table_name"
	AWSTableNameAttribute2 = "aws.table.name"

	// AWSXRayInProgressAttribute is the `in_progress` flag in an X-Ray segment
	AWSXRayInProgressAttribute = "aws.xray.inprogress"

	// AWSXRayXForwardedForAttribute is the `x_forwarded_for` flag in an X-Ray segment
	AWSXRayXForwardedForAttribute = "aws.xray.x_forwarded_for"

	// AWSXRayResourceARNAttribute is the `resource_arn` field in an X-Ray segment
	AWSXRayResourceARNAttribute = "aws.xray.resource_arn"

	// AWSXRayTracedAttribute is the `traced` field in an X-Ray subsegment
	AWSXRayTracedAttribute = "aws.xray.traced"

	// AWSXraySegmentAnnotationsAttribute is the attribute that
	// will be treated by the X-Ray exporter as the annotation keys.
	AWSXraySegmentAnnotationsAttribute = "aws.xray.annotations"

	// AWSXraySegmentMetadataAttributePrefix is the prefix of the attribute that
	// will be treated by the X-Ray exporter as metadata. The key of a metadata
	// will be AWSXraySegmentMetadataAttributePrefix + <metadata_key>.
	AWSXraySegmentMetadataAttributePrefix = "aws.xray.metadata."

	// AWSXrayRetriesAttribute is the `retries` field in an X-Ray (sub)segment.
	AWSXrayRetriesAttribute = "aws.xray.retries"

	// AWSXrayExceptionIDAttribute is the `id` field in an exception
	AWSXrayExceptionIDAttribute = "aws.xray.exception.id"
	// AWSXrayExceptionRemoteAttribute is the `remote` field in an exception
	AWSXrayExceptionRemoteAttribute = "aws.xray.exception.remote"
	// AWSXrayExceptionTruncatedAttribute is the `truncated` field in an exception
	AWSXrayExceptionTruncatedAttribute = "aws.xray.exception.truncated"
	// AWSXrayExceptionSkippedAttribute is the `skipped` field in an exception
	AWSXrayExceptionSkippedAttribute = "aws.xray.exception.skipped"
	// AWSXrayExceptionCauseAttribute is the `cause` field in an exception
	AWSXrayExceptionCauseAttribute = "aws.xray.exception.cause"
)
