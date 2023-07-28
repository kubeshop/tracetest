import {renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useGetTestByIdQuery} from '../Tracetest';

test('useGetTestByIdQuery', async () => {
  const testId = '34k23';
  fetchMock.mockResponse(JSON.stringify({id: testId}));
  const {result, waitForNextUpdate} = renderHook(() => useGetTestByIdQuery({testId}), {
    wrapper: ReduxWrapperProvider,
  });
  expect(result.current.isLoading).toBeTruthy();
  await waitForNextUpdate({timeout: 5000});
  expect(result.current.isLoading).toBeFalsy();
});
