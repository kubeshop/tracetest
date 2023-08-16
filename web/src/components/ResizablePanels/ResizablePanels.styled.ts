import {Button} from 'antd';
import styled, {createGlobalStyle, css} from 'styled-components';
import {withPulseAnimation} from '../PulseButton';

export const GlobalStyle = createGlobalStyle`
  .spaces-resize-handle {
    background-color: ${({theme}) => theme.color.borderLight};
    width: 3px !important;
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
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
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
