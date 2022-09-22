import styled from 'styled-components';

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
  overflow-y: ${({$shouldScroll}) => ($shouldScroll ? 'scroll' : 'hidden')};
  z-index: 2;
`;

export const SwitchContainer = styled.div`
  bottom: 163px;
  left: 16px;
  position: absolute;
  z-index: 9;
`;
