import {CSSProperties} from 'react';

export function doArrow(
  position: string,
  verticalAlign: string,
  horizontalAlign: string
): Record<string, CSSProperties> {
  const opositeSide: Record<string, string> = {top: 'bottom', bottom: 'top', right: 'left', left: 'right'};

  if (!position || position === 'custom') {
    return {};
  }

  const width = 16;
  const height = 12;
  const color = 'white';
  const isVertical = position === 'top' || position === 'bottom';
  const spaceFromSide = 10;

  return {
    '&::after': {
      content: "''",
      width: 0,
      height: 0,
      position: 'absolute',
      [isVertical ? 'borderLeft' : 'borderTop']: `${width / 2}px solid transparent`, // CSS Triangle width
      [isVertical ? 'borderRight' : 'borderBottom']: `${width / 2}px solid transparent`, // CSS Triangle width
      [`border${position[0].toUpperCase()}${position.substring(1)}`]: `${height}px solid ${color}`, // CSS Triangle height
      [isVertical ? opositeSide[horizontalAlign] : verticalAlign]: height + spaceFromSide, // space from side
      [opositeSide[position]]: -height + 2,
    },
  };
}
