import {useCallback} from 'react';

const useScrollTo = () => {
  const scrollTo = useCallback(({elementId, containerId}) => {
    const element = document.getElementById(elementId);
    const container = document.getElementById(containerId);
    const offsetTop = element?.offsetTop ?? 0;

    if (container) {
      container.scrollTop = offsetTop;
    }
  }, []);

  return scrollTo;
};

export default useScrollTo;
