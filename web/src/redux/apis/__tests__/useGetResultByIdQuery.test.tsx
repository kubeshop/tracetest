import {renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useGetRunByIdQuery} from '../TraceTest.api';

test('useGetResultByIdQuery', async () => {
  fetchMock.mockResponse(JSON.stringify({}));
  const {result, waitForNextUpdate} = renderHook(() => useGetRunByIdQuery({runId: '234', testId: '34k23'}), {
    wrapper: ReduxWrapperProvider,
  });
  expect(result.current.isLoading).toBeTruthy();
  await waitForNextUpdate({timeout: 10000});
  expect(result.current.isLoading).toBeFalsy();
});
