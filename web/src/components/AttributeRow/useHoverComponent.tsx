import {Dispatch, SetStateAction, useCallback} from 'react';

export const useHoverComponent = (
  index: number,
  isAnyHovered: number[],
  setIsAnyHovered: Dispatch<SetStateAction<number[]>>
) => {
  return {
    onMouseEnter: useCallback(() => {
      if (!isAnyHovered.includes(index)) {
        setIsAnyHovered(iah => [...iah, index]);
      }
    }, [index, setIsAnyHovered, isAnyHovered]),
    onMouseLeave: useCallback(() => {
      if (isAnyHovered.includes(index)) {
        setIsAnyHovered(iah => iah.filter(ah => ah !== index));
      }
    }, [index, setIsAnyHovered, isAnyHovered]),
  };
};
