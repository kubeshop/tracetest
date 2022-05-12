import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useRunTestMutation} from '../Test.api';

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

  let testResult;
  await act(async () => {
    const [runNewTest] = result.current;
    testResult = await runNewTest(`${testId}`).unwrap();
  });
  expect(testResult.testId).toBe(`${testId}`);
});
