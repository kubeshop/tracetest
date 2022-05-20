import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {TTestRun} from '../../../types/TestRun.types';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useRunTestMutation} from '../TraceTest.api';

test('useRunTestMutation', async () => {
  const testId = 22;
  fetchMock.mockResponse(
    JSON.stringify({
      assertionResult: undefined,
      assertionResultState: false,
      completedAt: '',
      createdAt: '',
      response: undefined,
      resultId: '',
      spanId: '',
      state: undefined,
      testId: `${testId}`,
      trace: undefined,
      traceId: '',
    })
  );
  const {result} = renderHook(() => useRunTestMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  let testResult: TTestRun;
  await act(async () => {
    const [runNewTest] = result.current;
    testResult = await runNewTest({testId: `${testId}`}).unwrap();
  });
  expect(testResult!.id).toBe(`${testId}`);
});
