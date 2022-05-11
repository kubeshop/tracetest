import {useEffect} from 'react';

interface usePollingProps {
  callback(): void;

  delay: number;
  isPolling: boolean;
}

const usePolling = ({callback, delay = 1000, isPolling}: usePollingProps): void => {
  useEffect(() => {
    let interval: any = null;
    interval = setInterval(() => {
      if (isPolling) {
        callback();
      } else {
        interval && clearInterval(interval);
      }
    }, delay);

    return () => {
      return interval && clearInterval(interval);
    };
  }, [delay, isPolling, callback]);
};

export default usePolling;
