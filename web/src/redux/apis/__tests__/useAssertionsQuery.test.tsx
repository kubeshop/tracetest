import {renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useGetAssertionsQuery} from '../Test.api';

test('useGetAssertionsQuery', async () => {
  fetchMock.mockResponse(JSON.stringify([{testId: '24'}]));
  const {result, waitForNextUpdate} = renderHook(() => useGetAssertionsQuery('231'), {
    wrapper: ReduxWrapperProvider,
  });
  expect(result.current.isLoading).toBeTruthy();
  await waitForNextUpdate({timeout: 5000});
  expect(result.current.isLoading).toBeFalsy();
  expect(result.current.data.length).toBe(1);
});
