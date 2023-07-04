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
  background-color: ${({theme}) => theme.color.white};
  overflow-y: scroll;
  z-index: 1;
  flex-basis: 50%;
  max-width: 50vw;
`;

export const SectionRight = styled(Section)`
  background-color: ${({theme}) => theme.color.white};
  border-left: 1px solid rgba(3, 24, 73, 0.1);
  overflow-y: scroll;
  z-index: 2;
  flex-basis: 50%;
  max-width: 50vw;
`;
