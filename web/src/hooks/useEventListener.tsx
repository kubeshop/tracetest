import {RefObject, useEffect, useRef} from 'react';
import {useIsomorphicLayoutEffect} from './useIsomorphicLayoutEffect';

export function useEventListener<K extends keyof HTMLElementEventMap, T extends HTMLElement = HTMLDivElement>(
  eventName: K,
  handler: (event: HTMLElementEventMap[K]) => void,
  element?: RefObject<T>
) {
  // Create a ref that stores handler
  const savedHandler = useRef(handler);

  useIsomorphicLayoutEffect(() => {
    savedHandler.current = handler;
  }, [handler]);

  useEffect(() => {
    // Define the listening target
    const targetElement: T | Window = element?.current || window;
    if (!(targetElement && targetElement.addEventListener)) {
      return;
    }

    // Create event listener that calls handler function stored in ref
    const eventListener: typeof handler = event => savedHandler.current(event);

    targetElement.addEventListener(eventName, eventListener as any);

    // Remove event listener on cleanup
    return () => {
      targetElement.removeEventListener(eventName, eventListener as any);
    };
  }, [eventName, element]);
}
