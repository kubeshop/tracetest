package tracedb

import (
	"context"
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
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	"go.opentelemetry.io/otel/trace"
)

type AwsCredentials struct {
	SecretAccessKey string
	AccessKey       string
}

type awsxrayDB struct {
	realTraceDB

	credentials *credentials.Credentials
	session     *session.Session
	region      string
	service     *xray.XRay
}

func NewAwsXRayDB(cfg *config.AWSXRayDataStoreConfig) (TraceDB, error) {
	sessionCredentials := credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	return &awsxrayDB{
		credentials: sessionCredentials,
		region:      cfg.Region,
	}, nil
}

func (db *awsxrayDB) Connect(ctx context.Context) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: db.credentials,
	})

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

func (db *awsxrayDB) TestConnection(ctx context.Context) connection.ConnectionTestResult {
	url := fmt.Sprintf("xray.%s.amazonaws.com:443", db.region)
	tester := connection.NewTester(
		connection.WithConnectivityTest(connection.ConnectivityStep(connection.ProtocolHTTP, url)),
		connection.WithPollingTest(connection.TracePollingTestStep(db)),
		connection.WithAuthenticationTest(connection.NewTestStep(func(ctx context.Context) (string, error) {
			_, err := db.GetTraceByID(ctx, id.NewRandGenerator().TraceID().String())
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

	trace, err := parseXRayTrace(traceID, res.Traces[0])

	return trace, err
}

func parseXRayTrace(traceID string, rawTrace *xray.Trace) (model.Trace, error) {
	if len(rawTrace.Segments) == 0 {
		return model.Trace{}, nil
	}

	segment := rawTrace.Segments[0]
	spans, err := parseSegmentToSpans([]byte(*segment.Document), traceID)

	if err != nil {
		return model.Trace{}, err
	}

	return model.NewTrace(traceID, spans), nil
}

func parseSegmentToSpans(rawSeg []byte, traceID string) ([]model.Span, error) {
	var seg Segment
	err := json.Unmarshal(rawSeg, &seg)
	if err != nil {
		return []model.Span{}, err
	}

	err = seg.Validate()
	if err != nil {
		return []model.Span{}, err
	}

	spans, err := segToSpans(seg, traceID, nil)
	return spans, err
}

func segToSpans(seg Segment, traceID string, parent *model.Span) ([]model.Span, error) {
	span, err := generateSpan(&seg, traceID, parent)
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

		// if seg.Cause != nil &&
		// 	populatedChildSpan.Status().Code() != ptrace.StatusCodeUnset {
		// 	// if seg.Cause is not nil, then one of the subsegments must contain a
		// 	// HTTP error code. Also, span.Status().Code() is already
		// 	// set to `StatusCodeUnknownError` by `addCause()` in
		// 	// `populateSpan()` above, so here we are just trying to figure out
		// 	// whether we can get an even more specific error code.

		// 	if span.Status().Code() == ptrace.StatusCodeError {
		// 		// update the error code to a possibly more specific code
		// 		span.Status().SetCode(populatedChildSpan.Status().Code())
		// 	}
		// }
	}

	return spans, nil
}

func generateSpan(seg *Segment, traceID string, parent *model.Span) (model.Span, error) {
	attributes := make(model.Attributes, 0)
	span := model.Span{
		Parent: parent,
		Name:   *seg.Name,
	}

	if parent != nil {
		attributes["parent_id"] = parent.ID.String()
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

	attributes.SetPointer(conventions.AttributeEnduserID, seg.User)
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

func addNamespace(seg *Segment, attributes model.Attributes) error {
	if seg.Namespace != nil {
		switch *seg.Namespace {
		case validAWSNamespace:
			attributes.SetPointer(AWSServiceAttribute, seg.Name)

		case validRemoteNamespace:
			// no op
		default:
			return fmt.Errorf("unexpected namespace: %s", *seg.Namespace)
		}
		return nil
	}

	return nil
}

func addHTTP(seg *Segment, attributes model.Attributes) {
	if seg.HTTP == nil {
		return
	}

	if req := seg.HTTP.Request; req != nil {
		attributes.SetPointer(conventions.AttributeHTTPMethod, req.Method)
		attributes.SetPointer(conventions.AttributeHTTPClientIP, req.ClientIP)
		attributes.SetPointer(conventions.AttributeHTTPUserAgent, req.UserAgent)
		attributes.SetPointer(conventions.AttributeHTTPURL, req.URL)

		if req.XForwardedFor != nil {
			attributes[AWSXRayXForwardedForAttribute] = strconv.FormatBool(*req.XForwardedFor)
		}
	}

	if resp := seg.HTTP.Response; resp != nil {
		if resp.Status != nil {
			attributes[conventions.AttributeHTTPStatusCode] = fmt.Sprintf("%v", *resp.Status)
		}

		switch val := resp.ContentLength.(type) {
		case string:
			attributes[conventions.AttributeHTTPResponseContentLength] = val
		case float64:
			lengthPointer := int64(val)
			attributes[conventions.AttributeHTTPResponseContentLength] = fmt.Sprintf("%v", lengthPointer)
		}
	}
}

func addAWSToSpan(aws *AWSData, attrs model.Attributes) {
	if aws != nil {
		attrs.SetPointer(AWSAccountAttribute, aws.AccountID)
		attrs.SetPointer(AWSOperationAttribute, aws.Operation)
		attrs.SetPointer(AWSRegionAttribute, aws.RemoteRegion)
		attrs.SetPointer(AWSRequestIDAttribute, aws.RequestID)
		attrs.SetPointer(AWSQueueURLAttribute, aws.QueueURL)
		attrs.SetPointer(AWSTableNameAttribute, aws.TableName)

		if aws.Retries != nil {
			attrs[AWSXrayRetriesAttribute] = fmt.Sprintf("%v", *aws.Retries)
		}
	}
}

func addSQLToSpan(sql *SQLData, attrs model.Attributes) error {
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
	attrs.SetPointer(conventions.AttributeDBSystem, sql.DatabaseType)
	attrs.SetPointer(conventions.AttributeDBStatement, sql.SanitizedQuery)
	attrs.SetPointer(conventions.AttributeDBUser, sql.User)
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
		// maxAge of 28 days.  AWS has a 30 day limit, let's be conservative rather than
		// hit the limit
		maxAge = 60 * 60 * 24 * 28

		// maxSkew allows for 5m of clock skew
		maxSkew = 60 * 5
	)

	var (
		content      = [traceIDLength]byte{}
		epochNow     = time.Now().Unix()
		traceIDBytes = traceID
		epoch        = int64(binary.BigEndian.Uint32(traceIDBytes[0:4]))
		b            = [4]byte{}
	)

	// If AWS traceID originally came from AWS, no problem.  However, if oc generated
	// the traceID, then the epoch may be outside the accepted AWS range of within the
	// past 30 days.
	//
	// In that case, we return invalid traceid error

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

func decodeXRayTraceID(traceID string) (trace.TraceID, error) {
	tid := [16]byte{}

	if len(traceID) < 35 {
		return tid, errors.New("traceID length is wrong")
	}
	traceIDtoBeDecoded := (traceID)[2:10] + (traceID)[11:]

	_, err := hex.Decode(tid[:], []byte(traceIDtoBeDecoded))
	return tid, err
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
	// CauseTypeExceptionID indicates that the type of the `cause`
	// field is a string
	CauseTypeExceptionID CauseType = iota + 1
	// CauseTypeObject indicates that the type of the `cause`
	// field is an object
	CauseTypeObject
)

// Segment schema is documented in xray-segmentdocument-schema-v1.0.0 listed
// on https://docs.aws.amazon.com/xray/latest/devguide/xray-api-segmentdocuments.html
type Segment struct {
	// Required fields for both segment and subsegments
	Name      *string  `json:"name"`
	ID        *string  `json:"id"`
	StartTime *float64 `json:"start_time"`

	// Segment-only optional fields
	Service     *ServiceData `json:"service,omitempty"`
	Origin      *string      `json:"origin,omitempty"`
	User        *string      `json:"user,omitempty"`
	ResourceARN *string      `json:"resource_arn,omitempty"`

	// Optional fields for both Segment and subsegments
	TraceID     *string                           `json:"trace_id,omitempty"`
	EndTime     *float64                          `json:"end_time,omitempty"`
	InProgress  *bool                             `json:"in_progress,omitempty"`
	HTTP        *HTTPData                         `json:"http,omitempty"`
	Fault       *bool                             `json:"fault,omitempty"`
	Error       *bool                             `json:"error,omitempty"`
	Throttle    *bool                             `json:"throttle,omitempty"`
	Cause       *CauseData                        `json:"cause,omitempty"`
	AWS         *AWSData                          `json:"aws,omitempty"`
	Annotations map[string]interface{}            `json:"annotations,omitempty"`
	Metadata    map[string]map[string]interface{} `json:"metadata,omitempty"`
	Subsegments []Segment                         `json:"subsegments,omitempty"`

	// (for both embedded and independent) subsegment-only (optional) fields.
	// Please refer to https://docs.aws.amazon.com/xray/latest/devguide/xray-api-segmentdocuments.html#api-segmentdocuments-subsegments
	// for more information on subsegment.
	Namespace    *string  `json:"namespace,omitempty"`
	ParentID     *string  `json:"parent_id,omitempty"`
	Type         *string  `json:"type,omitempty"`
	PrecursorIDs []string `json:"precursor_ids,omitempty"`
	Traced       *bool    `json:"traced,omitempty"`
	SQL          *SQLData `json:"sql,omitempty"`
}

// Validate checks whether the segment is valid or not
func (s *Segment) Validate() error {
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

// AWSData represents the aws resource that this segment
// originates from
type AWSData struct {
	// Segment-only
	Beanstalk *BeanstalkMetadata `json:"elastic_beanstalk,omitempty"`
	CWLogs    []LogGroupMetadata `json:"cloudwatch_logs,omitempty"`
	ECS       *ECSMetadata       `json:"ecs,omitempty"`
	EC2       *EC2Metadata       `json:"ec2,omitempty"`
	EKS       *EKSMetadata       `json:"eks,omitempty"`
	XRay      *XRayMetaData      `json:"xray,omitempty"`

	// For both segment and subsegments
	AccountID    *string `json:"account_id,omitempty"`
	Operation    *string `json:"operation,omitempty"`
	RemoteRegion *string `json:"region,omitempty"`
	RequestID    *string `json:"request_id,omitempty"`
	QueueURL     *string `json:"queue_url,omitempty"`
	TableName    *string `json:"table_name,omitempty"`
	Retries      *int64  `json:"retries,omitempty"`
}

// EC2Metadata represents the EC2 metadata field
type EC2Metadata struct {
	InstanceID       *string `json:"instance_id"`
	AvailabilityZone *string `json:"availability_zone"`
	InstanceSize     *string `json:"instance_size"`
	AmiID            *string `json:"ami_id"`
}

// ECSMetadata represents the ECS metadata field. All must be omitempty b/c they come from two different detectors:
// Docker and ECS, so it's possible one is present and not the other
type ECSMetadata struct {
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
type BeanstalkMetadata struct {
	Environment  *string `json:"environment_name"`
	VersionLabel *string `json:"version_label"`
	DeploymentID *int64  `json:"deployment_id"`
}

// EKSMetadata represents the EKS metadata field
type EKSMetadata struct {
	ClusterName *string `json:"cluster_name"`
	Pod         *string `json:"pod"`
	ContainerID *string `json:"container_id"`
}

// LogGroupMetadata represents a single CloudWatch Log Group
type LogGroupMetadata struct {
	LogGroup *string `json:"log_group"`
	Arn      *string `json:"arn,omitempty"`
}

// CauseData is the container that contains the `cause` field
type CauseData struct {
	Type CauseType `json:"-"`
	// it will contain one of ExceptionID or (WorkingDirectory, Paths, Exceptions)
	ExceptionID *string `json:"-"`

	CauseObject
}

type CauseObject struct {
	WorkingDirectory *string     `json:"working_directory,omitempty"`
	Paths            []string    `json:"paths,omitempty"`
	Exceptions       []Exception `json:"exceptions,omitempty"`
}

// UnmarshalJSON is the custom unmarshaller for the cause field
func (c *CauseData) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &c.CauseObject)
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
type Exception struct {
	ID        *string      `json:"id,omitempty"`
	Message   *string      `json:"message,omitempty"`
	Type      *string      `json:"type,omitempty"`
	Remote    *bool        `json:"remote,omitempty"`
	Truncated *int64       `json:"truncated,omitempty"`
	Skipped   *int64       `json:"skipped,omitempty"`
	Cause     *string      `json:"cause,omitempty"`
	Stack     []StackFrame `json:"stack,omitempty"`
}

// StackFrame represents a frame in the stack when an exception occurred
type StackFrame struct {
	Path  *string `json:"path,omitempty"`
	Line  *int    `json:"line,omitempty"`
	Label *string `json:"label,omitempty"`
}

// HTTPData provides the shape for unmarshalling request and response fields.
type HTTPData struct {
	Request  *RequestData  `json:"request,omitempty"`
	Response *ResponseData `json:"response,omitempty"`
}

// RequestData provides the shape for unmarshalling the request field.
type RequestData struct {
	// Available in segment
	XForwardedFor *bool `json:"x_forwarded_for,omitempty"`

	// Available in both segment and subsegments
	Method    *string `json:"method,omitempty"`
	URL       *string `json:"url,omitempty"`
	UserAgent *string `json:"user_agent,omitempty"`
	ClientIP  *string `json:"client_ip,omitempty"`
}

// ResponseData provides the shape for unmarshalling the response field.
type ResponseData struct {
	Status        *int64      `json:"status,omitempty"`
	ContentLength interface{} `json:"content_length,omitempty"`
}

// ECSData provides the shape for unmarshalling the ecs field.
type ECSData struct {
	Container *string `json:"container"`
}

// EC2Data provides the shape for unmarshalling the ec2 field.
type EC2Data struct {
	InstanceID       *string `json:"instance_id"`
	AvailabilityZone *string `json:"availability_zone"`
}

// ElasticBeanstalkData provides the shape for unmarshalling the elastic_beanstalk field.
type ElasticBeanstalkData struct {
	EnvironmentName *string `json:"environment_name"`
	VersionLabel    *string `json:"version_label"`
	DeploymentID    *int    `json:"deployment_id"`
}

// XRayMetaData provides the shape for unmarshalling the xray field
type XRayMetaData struct {
	SDK                 *string `json:"sdk,omitempty"`
	SDKVersion          *string `json:"sdk_version,omitempty"`
	AutoInstrumentation *bool   `json:"auto_instrumentation"`
}

// SQLData provides the shape for unmarshalling the sql field.
type SQLData struct {
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
type ServiceData struct {
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
