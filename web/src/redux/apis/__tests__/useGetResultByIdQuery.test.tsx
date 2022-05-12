import {renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useGetResultByIdQuery} from '../Test.api';

test('useGetResultByIdQuery', async () => {
  fetchMock.mockResponse(JSON.stringify({}));
  const {result, waitForNextUpdate} = renderHook(() => useGetResultByIdQuery({resultId: '234', testId: '34k23'}), {
    wrapper: ReduxWrapperProvider,
  });
  expect(result.current.isLoading).toBeTruthy();
  await waitForNextUpdate({timeout: 5000});
  expect(result.current.isLoading).toBeFalsy();
});
