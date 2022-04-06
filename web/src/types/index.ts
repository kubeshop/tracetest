export enum HTTP_METHOD {
  GET = 'GET',
  PUT = 'PUT',
  POST = 'POST',
  PATCH = 'PATCH',
  DELETE = 'DELETE',
  COPY = 'COPY',
  HEAD = 'HEAD',
  OPTIONS = 'OPTIONS',
  LINK = 'LINK',
  UNLINK = 'UNLINK',
  PURGE = 'PURGE',
  LOCK = 'LOCK',
  UNLOCK = 'UNLOCK',
  PROPFIND = 'PROPFIND',
  VIEW = 'VIEW',
}

export const enum LOCATION_NAME {
  RESOURCE_ATTRIBUTES = 'RESOURCE_ATTRIBUTES',
  INSTRUMENTATION_LIBRARY = 'INSTRUMENTATION_LIBRARY',
  SPAN = 'SPAN',
  SPAN_ATTRIBUTES = 'SPAN_ATTRIBUTES',
  SPAN_ID = 'SPAN_ID',
}

export const enum COMPARE_OPERATOR {
  EQUALS = 'EQUALS',
  LESSTHAN = 'LESSTHAN',
  GREATERTHAN = 'GREATERTHAN',
  NOTEQUALS = 'NOTEQUALS',
  GREATOREQUALS = 'GREATOREQUALS',
  LESSOREQUAL = 'LESSOREQUAL',
}

export type ISpanAttributes = Array<{
  key: string;
  value: {[key: string]: any};
}>;

export interface ISpan {
  traceId: string;
  spanId: string;
  name: string;
  kind: string;
  startTime: number;
  duration: number;
  attributes: ISpanAttributes;
  events: Array<{
    timeUnixNano: string;
    name: string;
    attributes: Array<{
      key: string;
      value: {[key: string]: any};
    }>;
  }>;
}

export interface IAttribute {
  id?: string;
  key: string;
  value: string;
  type: 'span' | 'resource';
}

export interface TestRunResult {
  resultId: string;
  testId: string;
  traceId: string;
  spanId: string;
  createdAt: string;
  completedAt: string;
  response: any;
  trace: ITrace;
}

export interface HTTPRequest {
  url: string;
  method: HTTP_METHOD;
  headers?: Array<{[key: string]: string}>;
  body?: string;
  auth?: any;
  proxy?: any;
  certificate?: any;
}

export interface Test {
  testId: string;
  name: string;
  description: string;
  serviceUnderTest: {
    id: string;
    request: HTTPRequest;
  };
  assertions: Array<Assertion>;
  lastTestResult: TestRunResult;
}

export interface ItemSelector {
  locationName: LOCATION_NAME;
  propertyName: string;
  value: string;
  valueType: string;
}

export interface SpanSelector {
  spanAssertionId?: string;
  locationName: LOCATION_NAME;
  propertyName: string;
  valueType: string;
  operator: COMPARE_OPERATOR;
  comparisonValue: string;
}
export interface Assertion {
  assertionId: string;
  selectors: Array<ItemSelector>;
  spanAssertions: Array<SpanSelector>;
}

export interface TestResult {
  id: string;
  successful: {
    id: string;
    operationName: string;
    duration: string;
    numOfSPans: number;
    attributes: Array<IAttribute>;
  };
  failed: {
    id: string;
    operationName: string;
    duration: string;
    numOfSPans: number;
    attributes: Array<IAttribute>;
  };
  timeStamp: Date;
}

export interface ITrace {
  resourceSpans: Array<ResourceSpan>;
  description: string;
}

export interface ResourceSpan {
  resource: Resource;
  instrumentationLibrarySpans: InstrumentationLibrarySpan[];
}

export interface InstrumentationLibrarySpan {
  instrumentationLibrary: InstrumentationLibrary;
  spans: Span[];
}

export interface Event {
  timeUnixNano: any;
  name: string;
  attributes: Attribute[];
}

export interface Span {
  traceId: string;
  spanId: string;
  parentSpanId: string;
  name: string;
  kind: number;
  startTimeUnixNano: any;
  endTimeUnixNano: any;
  attributes: Attribute[];
  status: any;
  events: Event[];
}

export interface Resource {
  attributes: Attribute[];
}

export interface InstrumentationLibrary {
  name: string;
  version: string;
}

export interface Attribute {
  key: string;
  value: {[key: string]: string};
}
export interface ITestResult {
  resultId: string;
  traceid: string;
  spanid: string;
  successful: {};
  failed: {};
  createdAt: string;
  completedAt: string;
}

export type TestId = string;

export type AssertionResult = {
  spanListAssertionResult: SpanAssertionResult[][];
  assertion: Assertion;
  spanCount: number;
};

export interface SpanAssertionResult extends SpanSelector {
  hasPassed: boolean;
  actualValue: string;
}

export type RecursivePartial<T> = {
  [P in keyof T]?: RecursivePartial<T[P]>;
};
