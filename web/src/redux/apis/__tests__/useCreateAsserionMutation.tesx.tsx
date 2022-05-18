import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {IAssertion} from '../../../types/Assertion.types';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useCreateAssertionMutation} from '../Test.api';

test('useCreateAssertionMutation', async () => {
  expect.assertions(1);
  const assertionId = 22;
  fetchMock.mockResponse(JSON.stringify({assertionId}));
  const {result} = renderHook(() => useCreateAssertionMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  let newAssertion: IAssertion;
  await act(async () => {
    const [createTest] = result.current;
    newAssertion = await createTest({
      testId: '22',
      assertion: {
        assertionId: '34',
        spanAssertions: [],
        selectors: [],
      },
    }).unwrap();
  });

  expect(newAssertion!.assertionId).toBe(assertionId);
});
