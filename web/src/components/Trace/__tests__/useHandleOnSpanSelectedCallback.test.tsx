import {act, renderHook} from '@testing-library/react-hooks';
import {TestingModels} from '../../../utils/TestingModels';
import {useHandleOnSpanSelectedCallback} from '../hooks/useHandleOnSpanSelectedCallback';

test('useHandleOnSpanSelectedCallback', () => {
  const addSelected = jest.fn();
  const selectedSpan = jest.fn();
  const {result} = renderHook(() =>
    useHandleOnSpanSelectedCallback(addSelected, TestingModels.testRunResult, selectedSpan)
  );

  act(() => {
    result.current(TestingModels.spanId);
  });

  expect(addSelected).toBeCalled();
  expect(addSelected).toBeCalledWith([{id: TestingModels.spanId}]);
  expect(selectedSpan).toBeCalled();
  expect(selectedSpan).toBeCalledWith(TestingModels.span);
});
