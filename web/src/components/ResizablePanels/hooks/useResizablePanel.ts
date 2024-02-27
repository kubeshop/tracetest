import {useCallback, useEffect, useState} from 'react';

export type TPanel = {
  name: string;
  isDefaultOpen?: boolean;
  minSize?(): number;
  maxSize?(): number;
  openSize?(): number;
  closeSize?(): number;
};

interface IProps {
  panel: TPanel;
}

export type TSize = Exclude<TPanel, 'component'> & {
  size: number;
  isOpen: boolean;
  isFullScreen: boolean;
  minSize(): number;
  maxSize(): number;
  openSize(): number;
  fullScreen(): number;
  closeSize(): number;
};

const defaultMaxSize = () => window.innerWidth;
const defaultCloseSize = () => 25;
const defaultMinSize = () => 350;
const defaultOpenSize = () => window.innerWidth / 3;

const useResizablePanel = ({
  panel: {
    minSize = defaultMinSize,
    maxSize = defaultMaxSize,
    closeSize = defaultCloseSize,
    openSize = defaultOpenSize,
    ...panel
  },
}: IProps) => {
  const [size, setSize] = useState<TSize>({
    ...panel,
    size: (panel.isDefaultOpen && openSize()) || minSize() || 0,
    isOpen: panel.isDefaultOpen || false,
    minSize,
    maxSize,
    openSize,
    closeSize,
    isFullScreen: false,
    fullScreen: () => window.innerWidth,
  });

  const onWindowResize = useCallback(() => {
    if (size.isOpen) {
      setSize({
        ...size,
        size: size.isFullScreen ? size.fullScreen() : size.openSize(),
      });
    } else {
      setSize({
        ...size,
        size: size.closeSize(),
      });
    }
  }, [size]);

  useEffect(() => {
    window.addEventListener('resize', onWindowResize);

    return () => {
      window.removeEventListener('resize', onWindowResize);
    };
  }, [onWindowResize]);

  const onStopResize = useCallback(
    newSize => {
      const currentSize = Number(newSize);
      const isOpen = currentSize > size.minSize();

      setSize({
        ...size,
        size: isOpen ? currentSize : size.closeSize(),
        isOpen,
      });
    },
    [size]
  );

  const onChange = useCallback(
    (width: number) => {
      setSize({
        ...size,
        size: width,
        isOpen: width > size.minSize(),
        isFullScreen: width === size.fullScreen(),
      });
    },
    [size]
  );

  return {onStopResize, size, onChange};
};

export default useResizablePanel;
