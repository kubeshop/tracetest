import {useEffect, useRef} from 'react';
import {toPercent} from 'utils/Common';
import * as S from './VerticalResizer.styled';
import {useVerticalResizer} from '../VerticalResizer.provider';

const VerticalResizer = () => {
  const {columnWidth, highlightPosition, initResizer} = useVerticalResizer();
  const rootRef = useRef<HTMLDivElement>(null);
  const draggerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const resizer = initResizer({rootRef});
    const handleMouseDown = resizer.initEventHandler();
    draggerRef.current?.addEventListener('mousedown', handleMouseDown);

    return () => {
      resizer.cleanup();
      draggerRef.current?.removeEventListener('mousedown', handleMouseDown);
    };
  }, [initResizer]);

  let draggerClass = '';
  let draggerStyle;

  if (highlightPosition !== null) {
    draggerClass = `${highlightPosition > columnWidth ? 'right' : 'left'}-dragging`;
    // Draw a highlight from the current dragged position to the original position
    const draggerLeft = toPercent(Math.min(columnWidth, highlightPosition));
    // Subtract 1px for draggerRight to deal with the right border being off
    // by 1px when dragging left
    const draggerRight = `calc(${toPercent(1 - Math.max(columnWidth, highlightPosition))} - 1px)`;
    draggerStyle = {left: draggerLeft, right: draggerRight};
  } else {
    draggerStyle = {left: toPercent(columnWidth)};
  }

  return (
    <S.VerticalResizer ref={rootRef}>
      <S.VerticalResizerDragger className={draggerClass} style={draggerStyle} ref={draggerRef} />
    </S.VerticalResizer>
  );
};

export default VerticalResizer;
