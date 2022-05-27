import {useCallback} from 'react';
import {useAppDispatch} from '../../redux/hooks';
import {revertDefinition} from '../../redux/slices/TestDefinition.slice';

export function useRevertDefinitionCallback(index: number): () => void {
  const dispatch = useAppDispatch();
  return useCallback(() => {
    return dispatch(revertDefinition({index}));
  }, [dispatch, index]);
}
