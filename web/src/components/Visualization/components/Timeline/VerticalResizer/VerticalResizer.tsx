import {useEffect, useRef, useState} from 'react';
import DraggableManager from 'utils/DragabbleManager';
import * as S from './VerticalResizer.styled';

const MIN_NAME_COLUMN_WIDTH = 0.15;
const MAX_NAME_COLUMN_WIDTH = 0.65;

interface IProps {
  nameColumnWidth: number;
  onNameColumnWidthChange(width: number): void;
}

const VerticalResizer = ({nameColumnWidth, onNameColumnWidthChange}: IProps) => {
  const rootRef = useRef<HTMLDivElement>(null);
  const resizerRef = useRef<HTMLDivElement>(null);
  const [dragPosition, setDragPosition] = useState<number | null>(null);

  useEffect(() => {
    const draggableManager = DraggableManager({
      getBounds: () => {
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
        setDragPosition(null);
        onNameColumnWidthChange(value);
      },
      onDragMove: ({value}) => {
        setDragPosition(value);
      },
    });

    const handleMouseDown = draggableManager.init();
    resizerRef.current?.addEventListener('mousedown', handleMouseDown);

    return () => {
      draggableManager.dispose();
      resizerRef.current?.removeEventListener('mousedown', handleMouseDown);
    };
  }, [onNameColumnWidthChange]);

  let isDraggingCls = '';
  let draggerStyle;

  if (dragPosition !== null) {
    isDraggingCls = 'dragging';
    // Draw a highlight from the current dragged position to the original position
    const draggerLeft = `${Math.min(nameColumnWidth, dragPosition) * 100}%`;
    // subtract 1px for draggerRight to deal with the right border being off
    // by 1px when dragging left
    const draggerRight = `calc(${(1 - Math.max(nameColumnWidth, dragPosition)) * 100}% - 1px)`;
    draggerStyle = {left: draggerLeft, right: draggerRight};
  } else {
    draggerStyle = {left: `${nameColumnWidth * 100}%`};
  }

  return (
    <S.VerticalResizer ref={rootRef}>
      <S.VerticalResizerDragger className={isDraggingCls} style={draggerStyle} ref={resizerRef} />
    </S.VerticalResizer>
  );
};

export default VerticalResizer;
