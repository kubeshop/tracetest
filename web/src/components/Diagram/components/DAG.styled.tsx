import {Button} from 'antd';
import styled, {css} from 'styled-components';

export const Container = styled.div<{$showAffected: boolean}>`
  position: relative;
  height: 100%;

  ${({$showAffected}) =>
    $showAffected &&
    css`
      .react-flow__node-TraceNode:not(.affected) > div {
        opacity: 0.5;
      }
    `}
`;

export const Controls = styled.div`
  background-color: #fbfbff;
  border-bottom-left-radius: 8px;
  position: absolute;
  right: 0;
  top: -14px;
  z-index: 9;
`;

export const ZoomButton = styled(Button)`
  color: rgba(3, 24, 73, 0.3);
  width: 24px;

  &:hover,
  &:focus {
    background-color: unset;
  }
`;
