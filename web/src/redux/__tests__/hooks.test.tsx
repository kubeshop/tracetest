import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {useAppDispatch} from '../hooks';
import {ReduxWrapperProvider} from '../ReduxWrapperProvider';

test('useAppDispatch', async () => {
  const testId = '34k23';
  fetchMock.mockResponse(JSON.stringify({testId}));
  const {result} = renderHook(() => useAppDispatch(), {
    wrapper: ReduxWrapperProvider,
  });
  act(() => {
    result.current({type: 'WHATEVER'});
  });
  expect(result.current.name).toBe('');
});
