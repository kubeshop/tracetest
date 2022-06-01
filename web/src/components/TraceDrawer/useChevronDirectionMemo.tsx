import {useMemo} from 'react';

export function useChevronDirectionMemo(height = 0, max = 0, min = 0) {
  return useMemo(() => {
    const isMax = (height || 0) > (max || 0) - 10;
    if (isMax) return true;
    const isMin = (height || 0) < (min || 0) + 10;
    if (isMin) {
      return false;
    }
    return false;
  }, [height, max, min]);
}
