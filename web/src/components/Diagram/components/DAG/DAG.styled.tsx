import {Button, Typography} from 'antd';
import styled, {css} from 'styled-components';

export const Container = styled.div<{$showAffected: boolean}>`
  position: relative;
  height: 100%;

  .react-flow__node-TraceNode.matched > div {
    box-shadow: 0 4px 8px #c9cedb;
  }

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

export const SelectorControls = styled(Controls)`
  left: 0;
`;

export const ZoomButton = styled(Button)`
  color: rgba(3, 24, 73, 0.3);
  width: 24px;

  &:hover,
  &:focus {
    background-color: unset;
  }
`;

export const ToggleButton = styled(ZoomButton)``;

export const FocusedText = styled(Typography.Text)`
  color: rgba(3, 24, 73, 0.3);
  margin-left: 8px;
`;
