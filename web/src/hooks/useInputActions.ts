import {RefObject, useEffect} from 'react';
import useOnClickOutside from './useOnClickOutside';

type Handler = (event: MouseEvent | KeyboardEvent) => void;

const useInputActions = <T extends HTMLElement = HTMLElement>(ref: RefObject<T>, handler: Handler) => {
  useOnClickOutside(ref, handler, 'mousedown');

  useEffect(() => {
    const element = ref.current;
    if (!element) return;

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        handler(event);
      }

      if (event.key === 'Enter') {
        handler(event);
      }
    };

    element.addEventListener('keydown', handleKeyDown);
    return () => {
      element.removeEventListener('keydown', handleKeyDown);
    };
  }, [handler, ref]);
};

export default useInputActions;
