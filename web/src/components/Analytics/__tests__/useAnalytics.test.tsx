import {renderHook} from '@testing-library/react-hooks';
import useAnalytics from '../useAnalytics';

describe('useAnalytics', () => {
  it('should render the hook', () => {
    const {result} = renderHook(() => useAnalytics());

    expect(result.current.isEnabled).toBeFalsy();
  });
});
