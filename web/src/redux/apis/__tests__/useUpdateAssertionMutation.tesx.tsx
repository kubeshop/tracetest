import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useUpdateAssertionMutation} from '../Test.api';

test('useUpdateAssertionMutation', async () => {
  const traceId = 23423;
  fetchMock.mockResponse(JSON.stringify({trace: traceId}));
  const {result} = renderHook(() => useUpdateAssertionMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  let testResult;
  await act(async () => {
    const [updateResult] = result.current;
    testResult = await updateResult({
      testId: `${traceId}`,
      assertionId: 'dfkjns',
      assertion: {selectors: []},
    }).unwrap();
  });
  expect(testResult).toStrictEqual({
    trace: traceId,
  });
});
