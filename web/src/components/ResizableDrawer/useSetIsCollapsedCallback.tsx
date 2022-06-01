import * as React from 'react';
import {useCallback} from 'react';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {DrawerState} from './ResizableDrawer';

export function useSetIsCollapsedCallback() {
  const {setDrawerState} = useAssertionForm();
  return useCallback<React.MouseEventHandler<HTMLDivElement>>(
    e => {
      e.stopPropagation();
      setDrawerState(c => (c === DrawerState.OPEN ? DrawerState.CLOSE : DrawerState.OPEN));
    },
    [setDrawerState]
  );
}
