import {useMemo} from 'react';

export function useReferenceMemo(visiblePortion: number) {
  return useMemo(() => {
    return {
      CLOSE: visiblePortion,
      INITIAL: window.outerHeight * 0.2,
      FORM: window.outerHeight * 0.4,
      OPEN: window.outerHeight * 0.5,
    };
  }, [visiblePortion]);
}
