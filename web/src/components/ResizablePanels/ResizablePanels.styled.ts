import {Button} from 'antd';
import styled, {createGlobalStyle, css} from 'styled-components';

export const GlobalStyle = createGlobalStyle`
  .panel-handle:has(button:hover) {
    * {
      cursor: pointer !important;
    }
  }
`;

export const ButtonContainer = styled.div<{$placement: 'left' | 'right'}>`
  position: relative;
  width: 4px;
  height: 100%;
  background: ${({theme}) => theme.color.white};

  :hover {
    background: ${({theme}) => theme.color.primaryHover};
  }

  ${({$placement}) =>
    ($placement === 'left' &&
      css`
        border-left: 1px solid ${({theme}) => theme.color.borderLight};

        :hover {
          border-left: 1px solid ${({theme}) => theme.color.primaryHover};
        }
      `) ||
    css`
      border-right: 1px solid ${({theme}) => theme.color.borderLight};

      :hover {
        border-right: 1px solid ${({theme}) => theme.color.primaryHover};
      }
    `}
`;

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
    ($isOpen &&
      css`
        > div {
          opacity: 1;
          pointer-events: auto;
        }
      `) ||
    css`
      :hover {
        background: ${({theme}) => theme.color.primaryHover};
      }
      cursor: pointer;
    `}
`;

export const ToggleButton = styled(Button).attrs({
  shape: 'circle',
  type: 'primary',
  onMouseDown: (e: React.MouseEvent) => e.stopPropagation(),
})`
  && {
    position: absolute;
    top: 60px;
    left: -14px;
    z-index: 1030;

    border: 3px solid ${({theme}) => theme.color.primaryLight};
    background-clip: padding-box;
    > span {
      font-size: ${({theme}) => theme.size.md};
    }
  }
`;
