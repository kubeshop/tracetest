import {useMemo} from 'react';

export function useReferenceMemo(visiblePortion: number) {
  return useMemo(() => {
    return {
      CLOSE: visiblePortion,
      // used when the page initially loads
      INITIAL: window.outerHeight * 0.2,
      // used when user open the assertion block form
      FORM: window.outerHeight * 0.4,
      // used to limit drawer height
      OPEN: window.outerHeight * 0.85,
    };
  }, [visiblePortion]);
}
