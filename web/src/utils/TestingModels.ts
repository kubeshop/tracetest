import {HTTP_METHOD} from '../constants/Common.constants';
import {CompareOperator} from '../constants/Operator.constants';
import {SemanticGroupNames} from '../constants/SemanticGroupNames.constants';
import {LOCATION_NAME} from '../constants/Span.constants';
import {SpanAttributeType} from '../constants/SpanAttribute.constants';
import {TestState} from '../constants/TestRunResult.constants';
import {IAssertion, IAssertionResult, IItemSelector, ISpanAssertionResult} from '../types/Assertion.types';
import {ISpan} from '../types/Span.types';
import {ITest} from '../types/Test.types';
import {ITestRunResult} from '../types/TestRunResult.types';
import {ITrace} from '../types/Trace.types';

interface IProps {
  mouseEvent: MouseEvent;
  assertionId: string;
  spanAssertionResult: ISpanAssertionResult;
  trace: ITrace;
  span: ISpan;
  test: ITest;
  assertion: IAssertion;
  assertionResult: IAssertionResult;
  testRunResult: ITestRunResult;
  resultId: string;
  testId: string;
  spanId: string;
}

const resultId = '0932u0e3';
const testId = '2309ei30';
const spanId = '4938u43';

const test: ITest = {
  assertions: [],
  description: '',
  lastTestResult: undefined,
  name: '',
  serviceUnderTest: {
    id: '',
    request: {
      auth: undefined,
      body: '',
      certificate: undefined,
      headers: undefined,
      method: HTTP_METHOD.GET,
      proxy: undefined,
      url: '',
    },
  },
  testId,
};

const span: ISpan = {
  attributeList: [],
  attributes: {},
  duration: 0,
  endTimeUnixNano: '',
  kind: '',
  name: '',
  parentSpanId: '',
  signature: [
    {
      locationName: LOCATION_NAME.SPAN,
      propertyName: '',
      value: '',
      valueType: '',
    },
  ],
  spanId,
  startTimeUnixNano: '',
  status: {code: ''},
  traceId: '',
  type: SemanticGroupNames.Http,
};

let selector: IItemSelector = {
  locationName: LOCATION_NAME.SPAN,
  propertyName: '',
  value: '',
  valueType: '',
};

const assertionResult: IAssertionResult = {
  assertion: {
    assertionId: '',
    selectors: [selector],
    spanAssertions: [],
  },
  spanListAssertionResult: [
    {
      resultList: [
        {
          actualValue: 'New request',
          comparisonValue: '',
          hasPassed: false,
          locationName: LOCATION_NAME.SPAN,
          operator: CompareOperator.EQUALS,
          propertyName: '',
          spanAssertionId: '',
          spanId: '',
          valueType: SpanAttributeType.intValue,
        },
      ],
      span,
    },
  ],
};

const testRunResult: ITestRunResult = {
  executionTime: 0,
  failedAssertionCount: 0,
  passedAssertionCount: 0,
  totalAssertionCount: 0,

  assertionResult: [{assertionId: '', spanAssertionResults: []}],
  assertionResultState: false,
  completedAt: '',
  createdAt: new Date().toISOString().toString(),
  response: undefined,
  resultId: '',
  spanId: '',
  state: TestState.CREATED,
  testId: '',
  trace: {
    description: '',
    resourceSpans: undefined,
    spans: [span],
  },
  traceId: '',
};

const spanAssertionResult: ISpanAssertionResult = {
  actualValue: 'New request',
  comparisonValue: '',
  hasPassed: false,
  locationName: LOCATION_NAME.INSTRUMENTATION_LIBRARY,
  operator: CompareOperator.EQUALS,
  propertyName: '',
  spanAssertionId: '',
  spanId: '',
  valueType: SpanAttributeType.doubleValue,
};

const assertionId = '234234';

const assertion: IAssertion = {
  assertionId,
  selectors: [
    {
      locationName: LOCATION_NAME.INSTRUMENTATION_LIBRARY,
      propertyName: '',
      value: '',
      valueType: '',
    },
  ],
  spanAssertions: [],
};
export const TestingModels: IProps = {
  mouseEvent: new MouseEvent('click', {
    bubbles: true,
    cancelable: true,
  }),
  assertionId,
  assertion,
  spanAssertionResult,
  assertionResult,
  test,
  testRunResult,
  span,
  resultId,
  testId,
  spanId,
  trace: {
    description: '',
    spans: [span],
  },
};

class NewTestingModel {
  test: ITest = test;
  assertionResult: IAssertionResult = assertionResult;
  testRunResult: ITestRunResult = testRunResult;
  span: ISpan = span;
  resultId: string = resultId;
  testId: string = testId;
  spanId: string = spanId;
  trace: ITrace = {
    description: '',
    spans: [span],
  };
}

export const TestingObj = new NewTestingModel();
