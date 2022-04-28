import { HTTP_METHOD } from "../../constants/Common.contants";
import { TAssertion } from "../Assertion/Assertion.types";
import { TTestRunResult } from "../TestRunResult/TestRunResult.types";

export type THTTPRequest = {
  url: string;
  method: HTTP_METHOD;
  headers?: Array<{[key: string]: string}>;
  body?: string;
  auth?: any;
  proxy?: any;
  certificate?: any;
};

export type TTest = {
  testId: string;
  name: string;
  description: string;
  serviceUnderTest: {
    id: string;
    request: THTTPRequest;
  };
  assertions: Array<TAssertion>;
  lastTestResult: TTestRunResult;
};
