import {MenuInfo} from 'rc-menu/lib/interface';
import {useCallback} from 'react';

type Return = (info: Partial<MenuInfo>) => void;

export function useOnDeleteCallback(onDelete: () => void): Return {
  return useCallback(
    ({domEvent}) => {
      domEvent?.stopPropagation();
      onDelete();
    },
    [onDelete]
  );
}
