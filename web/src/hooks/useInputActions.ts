import {RefObject, useCallback, useEffect} from 'react';
import useOnClickOutside from './useOnClickOutside';

type Handler = (event: MouseEvent | KeyboardEvent) => void;

const useInputActions = <T extends HTMLElement = HTMLElement>(ref: RefObject<T>, handler: Handler) => {
  useOnClickOutside(ref, handler);

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        handler(event);
      }

      if (event.key === 'Enter') {
        handler(event);
      }
    },
    [handler]
  );

  useEffect(() => {
    const element = ref.current;

    if (element) {
      element.addEventListener('keydown', handleKeyDown);

      return () => {
        element.removeEventListener('keydown', handleKeyDown);
      };
    }
  }, [handleKeyDown, ref]);
};

export default useInputActions;
