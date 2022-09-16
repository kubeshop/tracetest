import styled, {createGlobalStyle, css} from 'styled-components';

export const GlobalStyle = createGlobalStyle`
  .reflex-container.vertical > .reflex-splitter {
    border-left: 1px solid rgba(3, 24, 73, 0.1);
    border-right: 1px solid rgba(3, 24, 73, 0.1);
    position: relative;
  }
`;

export const ButtonContainer = styled.div`
  position: absolute;
  right: -12px;
  top: 48px;
`;

export const Content = styled.div<{$isOpen: boolean}>`
  background-color: ${({theme}) => theme.color.white};
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
  height: 100%;
  overflow: visible;
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
