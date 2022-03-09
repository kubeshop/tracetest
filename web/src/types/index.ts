export interface ISpan {
  traceId: string;
  spanId: string;
  name: string;
  kind: string;
  startTime: number;
  duration: number;
  attributes: Array<{
    key: string;
    value: {[key: string]: any};
  }>;
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

export interface Test {
  id: string;
  name: string;
  description: string;
  serviceUnderTest: {
    id?: string;
    url: string;
    auth?: string;
  };
  assertions: Array<Assertion>;
  repeats: number;
}
export interface Assertion {
  id: string;
  selector: string;
  comparable: string;
  operator: string;
  successful?: string;
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
  id: string;
  traceid: string;
  spanid: string;
  successful: {};
  failed: {};
  createdAt: string;
  completedAt: string;
}

export type TestId = string;
