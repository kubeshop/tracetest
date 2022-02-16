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
