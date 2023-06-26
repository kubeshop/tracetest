import {useCallback, useState} from 'react';
import {TPanel} from '../ResizablePanels';

interface IProps {
  panels: TPanel[];
}

export type TSize = Exclude<TPanel, 'component'> & {
  size: number;
  isOpen: boolean;
  minSize: number;
  maxSize: number;
};

type TSizeState = Record<string, TSize>;

const getInitialState = (panels: TPanel[]): TSizeState =>
  panels.reduce(
    (acc, {component, ...panel}) => ({
      ...acc,
      [panel.name]: {
        ...panel,
        size: panel.isDefaultOpen ? panel.maxSize : panel.minSize,
        isOpen: panel.isDefaultOpen,
      },
    }),
    {}
  );

const useResizablePanels = ({panels}: IProps) => {
  const [sizes, setSizes] = useState<TSizeState>(getInitialState(panels));

  const setSize = useCallback((size: TSize) => {
    setSizes(prev => ({
      ...prev,
      [size.name]: size,
    }));
  }, []);

  const toggle = useCallback(
    (size: TSize) => {
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
    },
    [setSize]
  );

  const onStopResize = useCallback(
    (size: TSize, newSize) => {
      const currentSize = Number(newSize);

      setSize({
        ...size,
        size: currentSize,
        isOpen: currentSize > size.minSize,
      });
    },
    [setSize]
  );

  return {toggle, onStopResize, sizes};
};

export default useResizablePanels;
