import {SemanticAttributes, SemanticResourceAttributes} from '@opentelemetry/semantic-conventions';

export const TraceTestAttributes = {
  TRACETEST_RESPONSE_BODY: 'tracetest.response.body',
  TRACETEST_RESPONSE_HEADERS: 'tracetest.response.headers',
  TRACETEST_RESPONSE_STATUS: 'tracetest.response.status',
  TRACETEST_SPAN_DURATION: 'tracetest.span.duration',
  TRACETEST_SPAN_TYPE: 'tracetest.span.type',
  NAME: 'name',
};

export const Attributes = {
  ...SemanticAttributes,
  ...SemanticResourceAttributes,
  ...TraceTestAttributes,
};

export * from '@opentelemetry/semantic-conventions';
