import styled, {createGlobalStyle} from 'styled-components';

export const GlobalStyle = createGlobalStyle`
  .spaces-resize-handle {
    background-color: ${({theme}) => theme.color.borderLight};
    width: 3px !important;
  }
`;

export const ButtonContainer = styled.div`
  position: absolute;
  right: -12px;
  top: 48px;
`;

export const SplitterContainer = styled.div``;
