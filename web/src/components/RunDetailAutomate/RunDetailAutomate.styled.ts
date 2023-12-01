import styled from 'styled-components';

export const Container = styled.div`
  display: flex;
  height: 100%;
  width: 100%;
`;

export const Section = styled.div`
  flex: 1;
  background-color: ${({theme}) => theme.color.white};
  overflow-y: scroll;
  flex-basis: 50%;
`;

export const SectionLeft = styled(Section)`
  z-index: 1;
`;

export const SectionRight = styled(Section)`
  border-left: 1px solid rgba(3, 24, 73, 0.1);
  z-index: 2;
  padding: 24px;
`;
