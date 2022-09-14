import styled, {css} from 'styled-components';

export const Aside = styled.div<{$isOpen: boolean}>`
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

export const AsideButtonContainer = styled.div`
  position: absolute;
  right: -12px;
  top: 48px;
`;

export const AsideContent = styled.div`
  height: 100%;
  overflow: hidden;
`;

export const Container = styled.div`
  display: flex;
  height: 100%;
  width: 100%;
`;

export const Section = styled.div`
  flex: 1;
`;

export const SectionLeft = styled(Section)`
  background-color: ${({theme}) => theme.color.background};
  z-index: 1;
`;

export const SectionRight = styled(Section)<{$shouldScroll: boolean}>`
  background-color: ${({theme}) => theme.color.white};
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
  overflow-y: ${({$shouldScroll}) => $shouldScroll ? 'scroll' : 'hidden'};
  z-index: 2;
`;
