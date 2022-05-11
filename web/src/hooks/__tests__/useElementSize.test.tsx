import {renderHook} from '@testing-library/react-hooks';
import {useElementSize} from '../useElementSize';

test('useElementSize', () => {
  const {result} = renderHook(() => useElementSize());
  expect(result.current[1].width).toBe(0);
  expect(result.current[1].height).toBe(0);
});
