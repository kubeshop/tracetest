import {Button, Typography} from 'antd';
import styled, {css} from 'styled-components';

export const Container = styled.div<{$showAffected: boolean}>`
  position: relative;
  height: 100%;

  .react-flow__attribution {
    visibility: hidden;
  }

  .react-flow__minimap {
    bottom: 0;
    background-color: ${({theme}) => theme.color.background};
    right: 0;
  }

  .react-flow__node-span > div.matched {
    box-shadow: ${({theme}) => `0 4px 8px ${theme.color.border}`};
  }

  ${({$showAffected}) =>
    $showAffected &&
    css`
      .react-flow__node-span > div:not(.affected) {
        opacity: 0.5;
      }
    `}
`;

export const Controls = styled.div`
  background-color: ${({theme}) => theme.color.background};
  border-bottom-left-radius: 8px;
  position: absolute;
  right: 0;
  top: -14px;
  z-index: 9;
`;

export const SelectorControls = styled(Controls)`
  left: 0;
`;

export const ZoomButton = styled(Button)<{$isActive?: boolean}>`
  color: ${({theme, $isActive}) => ($isActive ? theme.color.interactive : theme.color.textLight)};
  width: 24px;

  &:focus {
    background-color: unset;
    color: ${({theme, $isActive}) => ($isActive ? theme.color.interactive : theme.color.textLight)};
  }

  &:hover {
    background-color: unset;
    color: ${({theme}) => theme.color.text};
  }
`;

export const ToggleButton = styled(ZoomButton)`
  color: ${({theme}) => theme.color.primary};

  &:hover,
  &:focus {
    color: ${({theme}) => theme.color.primary};
  }
`;

export const FocusedText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.primary};
  margin-left: 8px;
`;
