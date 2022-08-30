import styled, {css} from 'styled-components';

export const ButtonContainer = styled.div`
  position: absolute;
  right: -12px;
  top: 48px;
`;

export const Container = styled.div<{$isOpen: boolean}>`
  background-color: ${({theme}) => theme.color.white};
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
  height: 100%;
  min-width: 30px;
  overflow: visible;
  position: relative;
  transition: width ease 0.2s, min-width ease 0.2s;
  width: 30px;
  z-index: 2;

  > div:first-child {
    opacity: 0;
  }

  ${({$isOpen}) =>
    $isOpen &&
    css`
      min-width: 270px;
      transition: width ease 0.2s 0.05s, min-width ease 0.2s 0.05s;
      width: 270px;

      > div:first-child {
        opacity: 1;
      }
    `}
`;

export const Content = styled.div`
  height: 100%;
  overflow: hidden;
`;
