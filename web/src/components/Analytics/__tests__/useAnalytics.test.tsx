import {renderHook} from '@testing-library/react-hooks';
import useAnalytics from '../useAnalytics';

test('useAnalytics', () => {
  const {result} = renderHook(() => useAnalytics());
  expect(result.current.isEnabled).toBeTruthy();
});
