type RefType = 'CHILD_OF';
export interface ISpan {
  traceID: string;
  spanID: string;
  operationName: string;
  references: Array<{
    refType: RefType;
    traceID: string;
    spanID: string;
  }>;
  startTime: number;
  duration: number;
  tags: Array<{
    key: string;
    type: string;
    value: unknown;
  }>;
  logs: Array<{
    timestamp: number;
    fields: Array<{
      key: string;
      type: string;
      value: unknown;
    }>;
  }>;
  processID: string;
  warnings: null;
}

interface IAttribute {
  id: string;
  key: string;
  value: string;
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
  operationName: string;
  duration: string;
  numOfSPans: number;
  attributes: Array<IAttribute>;
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
