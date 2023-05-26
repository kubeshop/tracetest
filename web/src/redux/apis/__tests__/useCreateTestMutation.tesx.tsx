import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {HTTP_METHOD} from '../../../constants/Common.constants';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useCreateTestMutation} from '../TraceTest.api';
import TestMock from '../../../models/__mocks__/Test.mock';

test('useCreateTestMutation', async () => {
  const test = TestMock.model();
  fetchMock.mockResponse(JSON.stringify(test));
  const {result} = renderHook(() => useCreateTestMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  await act(async () => {
    const [createTest] = result.current;
    const newTest = await createTest({
      name: 'New test',
      serviceUnderTest: {
        triggerType: 'http',
        http: {url: 'https://google.com', method: HTTP_METHOD.GET},
      },
    }).unwrap();
    expect(newTest!.id).toBe(test.id);
  });
});
