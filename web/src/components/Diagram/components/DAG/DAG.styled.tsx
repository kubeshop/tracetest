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
    border: ${({theme}) => `1px solid ${theme.color.text}`};
    box-shadow: ${({theme}) => `4px 4px 0px ${theme.color.text}`};
  }

  ${({$showAffected}) =>
    $showAffected &&
    css`
      .react-flow__node-span > div:not(.affected) {
        opacity: 0.5;
      }
    `}
`;

export const Controls = styled.div<{mode?: ControlsMode}>`
  background-color: ${({theme}) => theme.color.background};
  border-bottom-left-radius: 8px;
  position: ${({mode}) => (mode === 'timeline' ? 'unset' : 'absolute')};
  right: 0;
  top: -14px;
  z-index: 9;
`;

export type ControlsMode = 'timeline' | 'dag';

export const SelectorControls = styled(Controls)<{mode: ControlsMode}>`
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
