import styled, {createGlobalStyle, css} from 'styled-components';

export const GlobalStyle = createGlobalStyle`
  .spaces-space {
    border-left: 1px solid ${({theme}) => theme.color.borderLight};
  }

  .spaces-resize-handle.resize-right:before,
  .spaces-resize-handle.resize-left:before {
    cursor: col-resize;    
  }

  .splitter {
    .ant-tooltip-arrow-content:before,
    .ant-tooltip-inner {
      background: ${({theme}) => theme.color.primary};
      color: ${({theme}) => theme.color.white};
    }
  }
`;

export const SplitterContainer = styled.div`
  :hover {
    background: ${({theme}) => theme.color.primaryLight};
  }
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
        background: ${({theme}) => theme.color.primaryLight};
      }
      cursor: pointer;
    `}
`;
