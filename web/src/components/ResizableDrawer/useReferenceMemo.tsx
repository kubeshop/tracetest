import {useMemo} from 'react';
import {DrawerState} from './ResizableDrawer';

export function useReferenceMemo(visiblePortion: number) {
  return useMemo<Record<DrawerState, number>>(() => {
    return {
      CLOSE: visiblePortion,
      // value not used
      RESIZING: 0,
      // used when the page initially loads
      INITIAL: window.outerHeight * 0.2,
      // used when user open the assertion block form
      FORM: window.outerHeight * 0.4,
      // used when drawer is opened by clicking the header
      OPEN: window.outerHeight * 0.5,
      // used to limit drawer height
      MAX: window.outerHeight * 0.85,
    };
  }, [visiblePortion]);
}
