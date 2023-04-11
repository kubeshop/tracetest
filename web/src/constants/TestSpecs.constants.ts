import {IValues} from 'components/TestSpecForm/TestSpecForm';

export type TSnippet = Required<IValues>;

export const HTTP_SPANS_STATUS_CODE: TSnippet = {
  name: 'All HTTP Spans: Status  code is 200',
  selector: 'span[tracetest.span.type="http"]',
  assertions: [
    {
      left: 'attr:http.status_code',
      comparator: '=',
      right: '200',
    },
  ],
};

export const TRIGGER_SPAN_RESPONSE_TIME: TSnippet = {
  name: 'Trigger Span: Response time is less than 200ms',
  selector: 'span[tracetest.span.type="general" name="Tracetest trigger"]',
  assertions: [
    {
      left: 'attr:tracetest.span.duration',
      comparator: '<',
      right: '200ms',
    },
  ],
};

export const DB_SPANS_RESPONSE_TIME: TSnippet = {
  name: 'All Database Spans: Processing time is less than 100ms',
  selector: 'span[tracetest.span.type="database"]',
  assertions: [
    {
      left: 'attr:tracetest.span.duration',
      comparator: '<',
      right: '100ms',
    },
  ],
};

export const TRIGGER_SPAN_RESPONSE_BODY_CONTAINS: TSnippet = {
  name: 'Trigger Span: Response body contains "this string"',
  selector: 'span[tracetest.span.type="general" name="Tracetest trigger"]',
  assertions: [
    {
      left: 'attr:tracetest.response.body',
      comparator: 'contains',
      right: '"this string"',
    },
  ],
};

export const GRPC_SPANS_STATUS_CODE: TSnippet = {
  name: 'All gRPC Spans: Status is Ok',
  selector: 'span[tracetest.span.type="rpc" rpc.system="grpc"]',
  assertions: [
    {
      left: 'attr:grpc.status_code',
      comparator: '=',
      right: '0',
    },
  ],
};

export const DB_SPANS_QUALITY_DB_STATEMENT_PRESENT: TSnippet = {
  name: 'All Database Spans: db.statement should always be defined (QUALITY)',
  selector: 'span[tracetest.span.type="database"]',
  assertions: [
    {
      left: 'attr:db.system',
      comparator: '!=',
      right: '""',
    },
  ],
};

export const TEST_SPEC_SNIPPETS: TSnippet[] = [
  HTTP_SPANS_STATUS_CODE,
  GRPC_SPANS_STATUS_CODE,
  TRIGGER_SPAN_RESPONSE_TIME,
  TRIGGER_SPAN_RESPONSE_BODY_CONTAINS,
  DB_SPANS_RESPONSE_TIME,
  DB_SPANS_QUALITY_DB_STATEMENT_PRESENT,
];
