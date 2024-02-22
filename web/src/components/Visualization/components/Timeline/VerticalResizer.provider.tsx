import {INITIAL_NAME_COLUMN_WIDTH, MAX_NAME_COLUMN_WIDTH, MIN_NAME_COLUMN_WIDTH} from 'constants/Timeline.constants';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import DraggableManager from 'utils/DragabbleManager';

interface IContext {
  columnWidth: number;
  highlightPosition: number | null;
  initResizer({rootRef}: {rootRef: React.RefObject<HTMLDivElement>}): ReturnType<typeof DraggableManager>;
}

export const Context = createContext<IContext>({
  columnWidth: 0,
  highlightPosition: null,
  initResizer: () => DraggableManager({onGetBounds: () => ({clientXLeft: 0, width: 0})}),
});

interface IProps {
  children: React.ReactNode;
}

export const useVerticalResizer = () => useContext(Context);

const VerticalResizerProvider = ({children}: IProps) => {
  const [columnWidth, setColumnWidth] = useState(INITIAL_NAME_COLUMN_WIDTH);
  const [highlightPosition, setHighlightPosition] = useState<number | null>(null);

  const initResizer = useCallback(({rootRef}: {rootRef: React.RefObject<HTMLDivElement>}) => {
    const draggableManager = DraggableManager({
      onGetBounds: () => {
        if (!rootRef.current) {
          return {clientXLeft: 0, width: 0, maxValue: 0, minValue: 0};
        }

        const {left: clientXLeft, width} = rootRef.current.getBoundingClientRect();
        return {
          clientXLeft,
          width,
          maxValue: MAX_NAME_COLUMN_WIDTH,
          minValue: MIN_NAME_COLUMN_WIDTH,
        };
      },
      onDragEnd: ({value, resetBounds}) => {
        resetBounds();
        setHighlightPosition(null);
        setColumnWidth(value);
      },
      onDragMove: ({value}) => {
        setHighlightPosition(value);
      },
    });

    return draggableManager;
  }, []);

  const value = useMemo<IContext>(
    () => ({
      columnWidth,
      highlightPosition,
      initResizer,
    }),
    [columnWidth, highlightPosition, initResizer]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default VerticalResizerProvider;
