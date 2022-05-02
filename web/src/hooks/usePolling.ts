import {useEffect} from 'react';

interface usePollingProps {
  callback(): void;
  delay: number;
  isPolling: boolean;
}

const usePolling = ({callback, delay = 1000, isPolling}: usePollingProps) => {
  useEffect(() => {
    let interval: any = null;

    interval = setInterval(() => {
      if (isPolling) callback();
      else {
        interval && clearInterval(interval);
      }
    }, delay);

    return () => interval && clearInterval(interval);
  }, [delay, isPolling, callback]);
};

export default usePolling;
