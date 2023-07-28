import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useDeleteTestByIdMutation} from '../Tracetest';

test('useDeleteTestByIdMutation', async () => {
  const testId = 22;
  fetchMock.mockResponse(JSON.stringify(undefined));
  const {result} = renderHook(() => useDeleteTestByIdMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  let testResult;
  await act(async () => {
    const [deleteTestById] = result.current;
    testResult = await deleteTestById({testId: String(testId)}).unwrap();
  });
  expect(testResult).toBe(null);
});
