import {useEffect} from 'react';

interface usePollingProps {
  callback(): void;
  delay: number;
  isPolling: boolean;
}

const usePolling = ({callback, delay = 1000, isPolling}: usePollingProps): void => {
  useEffect(() => {
    let interval: any = null;

    console.log(1);
    interval = setInterval(() => {
      console.log(2);
      if (isPolling) {
        console.log(3);
        callback();
      } else {
        console.log(4);
        interval && clearInterval(interval);
      }
    }, delay);

    return () => {
      console.log(5);
      return interval && clearInterval(interval);
    };
  }, [delay, isPolling, callback]);
};

export default usePolling;
