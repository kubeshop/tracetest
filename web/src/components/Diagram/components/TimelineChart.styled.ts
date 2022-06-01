import styled, {css} from 'styled-components';

export const Container = styled.div<{$barHeight: number; $showAffected: boolean}>`
  overflow-y: scroll;
  
  .rect-svg {
    width: 100%;
    height: ${({$barHeight}) => `${$barHeight}px`};
    stroke: none;
    fill: none;

    &:hover {
      fill: rgb(229, 231, 239);
    }
  }

  .rect-svg-selected {
    fill: rgba(72, 88, 108, 0.1);
    stroke-width: 1px;
    stroke: rgba(72, 88, 108, 0.5);
  }

  .span-name {
    width: 180px;
    height ${({$barHeight}) => `${$barHeight}px`};
    fill: #000;
    font-size: 14px;
    pointer-events: none;
    alignment-baseline: middle;
    dominant-baseline: middle;
  }

  .span-duration {
    width: 100px;
    height: ${({$barHeight}) => `${$barHeight}px`};
    fill: #9AA3AB;
    font-size: 14px;
    pointer-events: none;
    alignment-baseline: middle;
    dominant-baseline: middle;
  }

  .duration-line {
    height: 10px;
    stroke: none;
    pointer-events: none;
  }

  .grey-line {
    width: 100%;
    height: 1px;
    stroke: none;
    fill: rgba(154, 163, 171, 0.7);
    pointer-events: none;
  }

  .node {
    height: ${({$barHeight}) => `${$barHeight}px`};
    cursor: pointer;
    pointer-events: bounding-box;
  }

  .checkpoint-mark {
    width: 1px;
    height: 10px;
    stroke: none;
    fill: rgb(213, 215, 224);
  }

  .cross-line {
    width: 100%;
    height: 1px;
    stroke: none;
    fill: rgb(213, 215, 224);
  }

  .duration-ms-text {
    text-anchor: end;
    fill: #9AA3AB;
  }

  .tick {
    fill: #9AA3AB;
    stroke: none;
    font-size: 14px;
  }

  ${({$showAffected}) =>
    $showAffected &&
    css`
      .rect-svg:not(.rect-svg-affected) ~ .span-name,
      .rect-svg:not(.rect-svg-affected) ~ .span-duration,
      .rect-svg:not(.rect-svg-affected) ~ .duration-line {
        opacity: 0.5;
      }
    `}
`;
