import { TAssertionResultList } from "../Assertion/Assertion.types";
import { TTrace } from "../Trace/Trace.types";
import { TestState } from "./TestRunResult.constants";

export type TAttribute = {
  id?: string;
  key: string;
  value: string;
  type: 'span' | 'resource';
}

export interface TTestResult {
  id: string;
  successful: {
    id: string;
    operationName: string;
    duration: string;
    numOfSPans: number;
    attributes: Array<TAttribute>;
  };
  failed: {
    id: string;
    operationName: string;
    duration: string;
    numOfSPans: number;
    attributes: Array<TAttribute>;
  };
  timeStamp: Date;
}

export type TTestRunResult = {
  resultId: string;
  testId: string;
  traceId: string;
  spanId: string;
  createdAt: string;
  completedAt: string;
  response: any;
  trace: TTrace;
  state: TestState;
  assertionResultState: boolean;
  assertionResult: TAssertionResultList;
};