import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useRunTestMutation} from '../Tracetest';

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
      id: `${testId}`,
      trace: undefined,
      traceId: '',
      allPassed: true,
      result: {},
    })
  );
  const {result} = renderHook(() => useRunTestMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  await act(async () => {
    const [runNewTest] = result.current;
    let testResult = await runNewTest({testId: `${testId}`}).unwrap();
    expect(testResult!.id).toBe(`${testId}`);
  });
});
