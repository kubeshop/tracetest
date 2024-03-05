import {useCallback, useRef, useState} from 'react';
import {ImperativePanelHandle} from 'react-resizable-panels';

export type TPanel = {
  isDefaultOpen?: boolean;
  minSize?(): number;
  maxSize?(): number;
  openSize?(): number;
  closeSize?(): number;
};

interface IProps {
  panel: TPanel;
}

export type TSize = {
  isOpen: boolean;
  isDefaultOpen: boolean;
  minSize(): number;
  maxSize(): number;
  openSize(): number;
  closeSize(): number;
};

const defaultMaxSize = () => 100;
const defaultCloseSize = () => 2;
const defaultMinSize = () => 15;
const defaultOpenSize = () => (window.innerWidth / 3 / window.innerWidth) * 100;

const useResizablePanel = ({
  panel: {
    minSize = defaultMinSize,
    maxSize = defaultMaxSize,
    closeSize = defaultCloseSize,
    openSize = defaultOpenSize,
    ...panel
  },
}: IProps) => {
  const ref = useRef<ImperativePanelHandle>(null);

  const [size, setSize] = useState<TSize>({
    isOpen: !!panel.isDefaultOpen,
    isDefaultOpen: !!panel.isDefaultOpen,
    minSize,
    maxSize,
    openSize,
    closeSize,
  });

  const onStopResize = useCallback(
    newSize => {
      const isOpen = Number(newSize) > size.minSize();
      setSize(prev => ({...prev, isOpen}));

      if (!isOpen) {
        ref.current?.resize(size.closeSize());
      }
    },
    [size]
  );

  const onChange = useCallback(
    (width: number) => {
      const isOpen = width > size.minSize();
      setSize(prev => ({...prev, isOpen}));

      ref.current?.resize(width);
    },
    [size]
  );

  return {onStopResize, size, onChange, ref};
};

export default useResizablePanel;
