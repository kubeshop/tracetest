import {useCallback, useState} from 'react';

const useHover = () => {
  const [isHovering, setIsHovering] = useState(false);

  const onMouseEnter = useCallback(() => {
    setIsHovering(true);
  }, []);

  const onMouseLeave = useCallback(() => {
    setIsHovering(false);
  }, []);

  return {isHovering, onMouseEnter, onMouseLeave};
};

export default useHover;
