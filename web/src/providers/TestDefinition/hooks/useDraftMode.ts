import {useEffect} from 'react';

const useDraftMode = (isDraftMode: boolean) => {
  useEffect(() => {
    if (isDraftMode) window.onbeforeunload = () => true;
    else window.onbeforeunload = null;
  }, [isDraftMode]);

  useEffect(() => {
    return () => {
      window.onbeforeunload = null;
    };
  }, []);
};

export default useDraftMode;
