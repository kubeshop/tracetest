import {renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useGetTestListQuery} from '../Tracetest';

test('useGetTestListQuery', async () => {
  fetchMock.mockResponse(JSON.stringify([{testId: '24'}]));
  const {result, waitForNextUpdate} = renderHook(() => useGetTestListQuery({}), {
    wrapper: ReduxWrapperProvider,
  });
  expect(result.current.isLoading).toBeTruthy();
  await waitForNextUpdate({timeout: 5000});
  expect(result.current.isLoading).toBeFalsy();
});
