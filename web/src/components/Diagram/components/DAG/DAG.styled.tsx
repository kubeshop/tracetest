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

export const Panel = styled.div`
  background-color: ${({theme}) => theme.color.background};
  border-radius: 8px;
  padding: 0 8px;
  position: absolute;
  top: 16px;
  z-index: 9;
`;

export const DAGActionsPanel = styled(Panel)`
  right: 16px;
`;

export const NavigatorPanel = styled(Panel)`
  left: 16px;
`;

export const ActionButton = styled(Button)<{$isActive?: boolean}>`
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

export const ToggleButton = styled(ActionButton)`
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
