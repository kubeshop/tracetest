import {useCallback, useState} from 'react';

export type TPanel = {
  name: string;
  isDefaultOpen?: boolean;
  minSize?: number;
  maxSize?: number;
};

interface IProps {
  panel: TPanel;
}

export type TSize = Exclude<TPanel, 'component'> & {
  size: number;
  isOpen: boolean;
  minSize: number;
  maxSize: number;
};

const useResizablePanel = ({panel}: IProps) => {
  const [size, setSize] = useState<TSize>({
    ...panel,
    size: (panel.isDefaultOpen && panel.maxSize) || panel.minSize || 0,
    isOpen: panel.isDefaultOpen || false,
    minSize: panel.minSize || 0,
    maxSize: panel.maxSize || 0,
  });

  const toggle = useCallback(() => {
    if (size.size <= size.minSize) {
      const currentSize = size.maxSize;
      setSize({
        ...size,
        size: currentSize,
        isOpen: currentSize > size.minSize,
      });
      return;
    }
    setSize({
      ...size,
      size: size.minSize || 0,
      isOpen: false,
    });
  }, [size]);

  const onStopResize = useCallback(
    newSize => {
      const currentSize = Number(newSize);

      setSize({
        ...size,
        size: currentSize,
        isOpen: currentSize > size.minSize,
      });
    },
    [size]
  );

  return {toggle, onStopResize, size};
};

export default useResizablePanel;
