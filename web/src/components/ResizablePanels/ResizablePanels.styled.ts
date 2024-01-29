import {Button} from 'antd';
import styled, {createGlobalStyle, css} from 'styled-components';
import {withPulseAnimation} from '../PulseButton';

export const GlobalStyle = createGlobalStyle`
  .spaces-resize-handle {
    border-left: 1px solid ${({theme}) => theme.color.borderLight};
    z-index: 10;
  }

  .splitter {
    .ant-tooltip-arrow-content:before,
    .ant-tooltip-inner {
      background: ${({theme}) => theme.color.primary};
      color: ${({theme}) => theme.color.white};
    }
  }
`;

export const ButtonContainer = styled.div`
  position: absolute;
  right: -15px;
  top: 48px;
  z-index: 100;
`;

export const SplitterButton = styled(Button)<{$isPulsing: boolean}>`
  && {
    border: 3px solid ${({theme}) => theme.color.primaryLight};
    background-clip: padding-box;
    > span {
      font-size: ${({theme}) => theme.size.md};
    }
  }

  ${({theme, $isPulsing}) => $isPulsing && withPulseAnimation(theme)}
`;

export const SplitterContainer = styled.div``;

export const PanelContainer = styled.div<{$isOpen: boolean}>`
  background-color: ${({theme}) => theme.color.white};
  height: 100%;
  overflow: visible;
  overflow-y: scroll;
  position: relative;

  > div {
    opacity: 0;
    pointer-events: none;
  }

  ${({$isOpen}) =>
    $isOpen &&
    css`
      > div {
        opacity: 1;
        pointer-events: auto;
      }
    `}
`;
