import {act, renderHook} from '@testing-library/react-hooks';
import {useDoubleClick} from '../useDoubleClick';

test('useDoubleClick', () => {
  const doubleClick = jest.fn();
  const click = jest.fn();
  const {result} = renderHook(() => useDoubleClick(doubleClick, click));

  act(() => {
    result.current({detail: 28});
  });

  expect(doubleClick).toBeCalled();
});
