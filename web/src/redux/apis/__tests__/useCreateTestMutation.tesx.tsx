import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useCreateTestMutation} from '../Test.api';

test('useCreateTestMutation', async () => {
  const testId = 22;
  fetchMock.mockResponse(JSON.stringify({testId}));
  const {result} = renderHook(() => useCreateTestMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  let newTest;
  await act(async () => {
    const [createTest] = result.current;
    newTest = await createTest({
      name: 'New test',
      serviceUnderTest: {
        request: {url: 'https://google.com', method: 'GET'},
      },
    }).unwrap();
  });
  expect(newTest.testId).toBe(testId);
});
