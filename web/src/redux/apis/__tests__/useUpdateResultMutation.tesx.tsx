import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useUpdateResultMutation} from '../Test.api';

test('useUpdateResultMutation', async () => {
  const testId = 22;
  fetchMock.mockResponse(JSON.stringify({trace: 23423}));
  const {result} = renderHook(() => useUpdateResultMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  let testResult;
  await act(async () => {
    const [updateResult] = result.current;
    testResult = await updateResult({
      testId: `${testId}`,
      resultId: 'dfkjns',
      assertionResult: {assertionResult: [], assertionResultState: true},
    }).unwrap();
  });
  expect(testResult).toStrictEqual({
    trace: {
      description: '',
      spans: [],
      executionTime: 0,
      failedAssertionCount: 0,
      passedAssertionCount: 0,
      totalAssertionCount: 0,
    },
  });
});
