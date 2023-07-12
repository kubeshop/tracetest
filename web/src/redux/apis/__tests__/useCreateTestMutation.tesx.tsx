import {act, renderHook} from '@testing-library/react-hooks';
import fetchMock from 'jest-fetch-mock';
import {HTTP_METHOD} from 'constants/Common.constants';
import TestMock from 'models/__mocks__/Test.mock';
import {ReduxWrapperProvider} from '../../ReduxWrapperProvider';
import {useCreateTestMutation} from '../TraceTest.api';

test('useCreateTestMutation', async () => {
  const test = TestMock.raw();
  fetchMock.mockResponse(JSON.stringify(test));
  const {result} = renderHook(() => useCreateTestMutation(), {
    wrapper: ReduxWrapperProvider,
  });

  await act(async () => {
    const [createTest] = result.current;
    const newTest = await createTest({
      type: 'Test',
      spec: {
        name: 'New test',
        trigger: {
          type: 'http',
          httpRequest: {url: 'https://google.com', method: HTTP_METHOD.GET},
        },
      },
    }).unwrap();
    expect(newTest?.id).toBe(test.spec?.id);
  });
});
