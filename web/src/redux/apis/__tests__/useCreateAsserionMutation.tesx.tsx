import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useCreateAssertionMutation} from '../Test.api';

test('useCreateAssertionMutation', async () => {
  const testId = 22;
  fetchMock.mockResponse(JSON.stringify({testId}));
  const {result} = renderHook(() => useCreateAssertionMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  let newTest;
  await act(async () => {
    const [createTest] = result.current;
    newTest = await createTest({
      testId: `${testId}`,
      assertion: {
        assertionId: '34',
        spanAssertions: [],
        selectors: [],
      },
    }).unwrap();
  });
  expect(newTest.testId).toBe(testId);
});
