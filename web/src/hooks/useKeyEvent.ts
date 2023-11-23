import {useCallback, useEffect} from 'react';

type TEvent = 'keydown' | 'keypress' | 'keyup';

export enum Keys {
  Enter = 'Enter',
  Escape = 'Escape',
}

const useKeyEvent = (targetKey: string[], callback: () => void, events: TEvent[] = ['keydown']) => {
  const onEvent = useCallback(
    (e: KeyboardEvent) => {
      if (targetKey.includes(e.key)) {
        e.stopPropagation();
        callback();
      }
    },
    [callback, targetKey]
  );

  useEffect(() => {
    events.forEach(event => {
      window.addEventListener(event, onEvent);
    });

    return () => {
      events.forEach(event => {
        window.removeEventListener(event, onEvent);
      });
    };
  }, [events, onEvent]);
};

export default useKeyEvent;
