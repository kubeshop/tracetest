import {useEffect} from 'react';

const useBlockNavigation = (isBlocking: boolean) => {
  useEffect(() => {
    if (isBlocking) window.onbeforeunload = () => true;
    else window.onbeforeunload = null;
  }, [isBlocking]);

  useEffect(() => {
    return () => {
      window.onbeforeunload = null;
    };
  }, []);
};

export default useBlockNavigation;
