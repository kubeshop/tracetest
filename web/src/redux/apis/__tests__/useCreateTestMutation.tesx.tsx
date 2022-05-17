import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {HTTP_METHOD} from '../../../constants/Common.constants';
import { ITest } from '../../../types/Test.types';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useCreateTestMutation} from '../Test.api';

test('useCreateTestMutation', async () => {
  const testId = 22;
  fetchMock.mockResponse(JSON.stringify({testId}));
  const {result} = renderHook(() => useCreateTestMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  let newTest: ITest;
  await act(async () => {
    const [createTest] = result.current;
    newTest = await createTest({
      name: 'New test',
      serviceUnderTest: {
        request: {url: 'https://google.com', method: HTTP_METHOD.GET},
      },
    }).unwrap();
  });
  expect(newTest!.testId).toBe(testId);
});
