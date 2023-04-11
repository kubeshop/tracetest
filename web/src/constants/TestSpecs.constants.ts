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

export const TEST_SPEC_SNIPPETS: TSnippet[] = [
  HTTP_SPANS_STATUS_CODE,
  TRIGGER_SPAN_RESPONSE_TIME,
  DB_SPANS_RESPONSE_TIME,
];
