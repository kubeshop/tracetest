import {SemanticAttributes, SemanticResourceAttributes} from '@opentelemetry/semantic-conventions';

export const TraceTestAttributes = {
  NAME: 'tracetest.span.name',
  KIND: 'tracetest.span.kind',
  TRACETEST_SPAN_TYPE: 'tracetest.span.type',
  TRACETEST_SPAN_DURATION: 'tracetest.span.duration',
  TRACETEST_RESPONSE_STATUS: 'tracetest.response.status',
  TRACETEST_RESPONSE_BODY: 'tracetest.response.body',
  TRACETEST_RESPONSE_HEADERS: 'tracetest.response.headers',
  TRACETEST_SELECTED_SPANS_COUNT: 'tracetest.selected_spans.count',
};

export const Attributes: Record<string, string> = {
  ...SemanticAttributes,
  ...SemanticResourceAttributes,
  ...TraceTestAttributes,
  HTTP_REQUEST_HEADER: 'http.request.header.',
  HTTP_RESPONSE_HEADER: 'http.response.header',
};

export * from '@opentelemetry/semantic-conventions';
