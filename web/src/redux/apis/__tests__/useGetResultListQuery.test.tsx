import {renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useGetResultListQuery} from '../Test.api';

test('useGetResultListQuery', async () => {
  fetchMock.mockResponse(JSON.stringify([]));
  const {result, waitForNextUpdate} = renderHook(() => useGetResultListQuery({testId: '3321'}), {
    wrapper: ReduxWrapperProvider,
  });
  expect(result.current.isLoading).toBeTruthy();
  await waitForNextUpdate({timeout: 5000});
  expect(result.current.isLoading).toBeFalsy();
});
