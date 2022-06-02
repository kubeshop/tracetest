import * as React from 'react';
import {useCallback} from 'react';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {DrawerState} from './ResizableDrawer';

export function useSetIsCollapsedCallbackDirection(input: boolean) {
  const {setDrawerState} = useAssertionForm();
  return useCallback<React.MouseEventHandler<HTMLDivElement>>(
    e => {
      e.stopPropagation();
      setDrawerState(input ? DrawerState.CLOSE : DrawerState.OPEN);
    },
    [setDrawerState, input]
  );
}

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
