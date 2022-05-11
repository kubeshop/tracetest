import {renderHook} from '@testing-library/react-hooks';
import {GuidedTours} from '../services/GuidedTour.service';
import useGuidedTour from './useGuidedTour';

test('useGuidedTour', () => {
  const {result} = renderHook(() => useGuidedTour(GuidedTours.Home));
  expect(result.current.isOpen).toBe(false);
});
