import {MAX_NAME_COLUMN_WIDTH, MIN_NAME_COLUMN_WIDTH} from 'constants/Timeline.constants';
import {useEffect, useRef, useState} from 'react';
import {toPercent} from 'utils/Common';
import DraggableManager from 'utils/DragabbleManager';
import * as S from './VerticalResizer.styled';

interface IProps {
  nameColumnWidth: number;
  onNameColumnWidthChange(width: number): void;
}

const VerticalResizer = ({nameColumnWidth, onNameColumnWidthChange}: IProps) => {
  const rootRef = useRef<HTMLDivElement>(null);
  const draggerRef = useRef<HTMLDivElement>(null);
  const [dragPosition, setDragPosition] = useState<number | null>(null);

  useEffect(() => {
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
        setDragPosition(null);
        onNameColumnWidthChange(value);
      },
      onDragMove: ({value}) => {
        setDragPosition(value);
      },
    });

    const handleMouseDown = draggableManager.init();
    draggerRef.current?.addEventListener('mousedown', handleMouseDown);

    return () => {
      draggableManager.dispose();
      draggerRef.current?.removeEventListener('mousedown', handleMouseDown);
    };
  }, [onNameColumnWidthChange]);

  let draggerClass = '';
  let draggerStyle;

  if (dragPosition !== null) {
    draggerClass = `${dragPosition > nameColumnWidth ? 'right' : 'left'}-dragging`;
    // Draw a highlight from the current dragged position to the original position
    const draggerLeft = toPercent(Math.min(nameColumnWidth, dragPosition));
    // Subtract 1px for draggerRight to deal with the right border being off
    // by 1px when dragging left
    const draggerRight = `calc(${toPercent(1 - Math.max(nameColumnWidth, dragPosition))} - 1px)`;
    draggerStyle = {left: draggerLeft, right: draggerRight};
  } else {
    draggerStyle = {left: toPercent(nameColumnWidth)};
  }

  return (
    <S.VerticalResizer ref={rootRef}>
      <S.VerticalResizerDragger className={draggerClass} style={draggerStyle} ref={draggerRef} />
    </S.VerticalResizer>
  );
};

export default VerticalResizer;
