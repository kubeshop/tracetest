import {useCallback} from 'react';
import {useStore} from 'react-redux';

export function useGetIsSelectedCallback() {
  const store = useStore();
  return useCallback(
    (spanId: string): boolean => {
      const state = store.getState();
      const {selectedElements} = state;
      const found = selectedElements ? selectedElements.find(({id}: {id: string}) => id === spanId) : undefined;
      return Boolean(found);
    },
    [store]
  );
}
