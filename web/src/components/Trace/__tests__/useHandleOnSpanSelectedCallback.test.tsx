import {act, renderHook} from '@testing-library/react-hooks';
import SpanMock from '../../../models/__mocks__/Span.mock';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';
import {useHandleOnSpanSelectedCallback} from '../hooks/useHandleOnSpanSelectedCallback';

test('useHandleOnSpanSelectedCallback', () => {
  const spanId = '12345';
  const addSelected = jest.fn();
  const selectedSpan = jest.fn();
  const {result} = renderHook(() => useHandleOnSpanSelectedCallback(addSelected, TestRunMock.model(), selectedSpan));

  act(() => {
    result.current(spanId);
  });

  expect(addSelected).toBeCalled();
  expect(addSelected).toBeCalledWith([{id: spanId}]);
  expect(selectedSpan).toBeCalled();
});
