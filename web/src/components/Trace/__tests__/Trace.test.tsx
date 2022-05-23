import {render, waitFor} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import Trace from '../Trace';
import {SemanticGroupNames} from '../../../constants/SemanticGroupNames.constants';
import {TSpan} from '../../../types/Span.types';
import {TTestRun} from '../../../types/TestRun.types';
import {TTest} from '../../../types/Test.types';
import {HTTP_METHOD} from '../../../constants/Common.constants';

const spanId = '4938u43';

const span: TSpan = {
  attributeList: [],
  attributes: {},
  duration: 0,
  // endTimeUnixNano: '',
  // kind: '',
  name: '',
  // parentSpanId: '',
  signature: [
    {
      // locationName: SemanticGroupNames.Http,
      // propertyName: '',
      key: '',
      value: '',
      // valueType: '',
    },
  ],
  id: spanId,
  // startTimeUnixNano: '',
  // status: {code: ''},
  // traceId: '',
  type: SemanticGroupNames.Http,
  children: [],
  endTime: '',
  startTime: '',
};

const testRunResult: TTestRun = {
  executionTime: 0,
  failedAssertionCount: 0,
  passedAssertionCount: 0,
  totalAssertionCount: 0,

  // assertionResult: [{assertionId: '', spanAssertionResults: []}],
  // assertionResultState: false,
  completedAt: '',
  createdAt: new Date().toISOString().toString(),
  response: undefined,
  // resultId: '',
  spanId: '',
  state: 'CREATED',
  // testId: '',
  trace: {
    traceId: '',
    // traceId,
    // description: '',
    // resourceSpans: undefined,
    spans: [span],
  },
  result: {
    allPassed: true,
    resultList: [],
    results: undefined,
  },
  traceId: '',
};
const ttest: TTest = {
  // assertions: [],
  description: '',
  // lastTestResult: undefined,
  name: '',
  serviceUnderTest: {
    // id: '',
    request: {
      auth: undefined,
      body: '',
      // certificate: undefined,
      headers: undefined,
      method: HTTP_METHOD.GET,
      // proxy: undefined,
      url: '',
    },
  },
  definition: undefined,
  id: '',
  referenceTestRun: undefined,

  // testId,
};

test('Trace', async () => {
  const {getByText} = render(
    <ReactFlowProvider>
      <ReduxWrapperProvider>
        <div style={{width: 600, height: 600}}>
          <Trace minHeight="300px" run={testRunResult} test={ttest} visiblePortion={100} displayError={false} />
        </div>
      </ReduxWrapperProvider>
    </ReactFlowProvider>
  );

  await waitFor(() => getByText('HTTP'));
});
