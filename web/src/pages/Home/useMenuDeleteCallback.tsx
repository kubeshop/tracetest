import {useCallback} from 'react';
import {MenuInfo} from 'rc-menu/lib/interface';
import {ITest} from '../../types/Test.types';

export function useMenuDeleteCallback(
  deleteTestMutation: (testId: string) => void
): (iTest: ITest) => (menu: MenuInfo) => void {
  return useCallback(
    (j: ITest) => async (e: MenuInfo) => {
      e.domEvent.stopPropagation();
      await deleteTestMutation(j.testId);
    },
    [deleteTestMutation]
  );
}
