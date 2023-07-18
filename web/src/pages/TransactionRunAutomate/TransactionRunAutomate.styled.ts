import styled from 'styled-components';

export const Container = styled.div`
  background: ${({theme}) => theme.color.white};
  display: flex;
  flex: auto;
  min-height: 0;
  width: 100%;
`;

export const Section = styled.div`
  background-color: ${({theme}) => theme.color.white};
  flex-basis: 50%;
  overflow-y: auto;
`;

export const SectionLeft = styled(Section)`
  z-index: 1;
`;

export const SectionRight = styled(Section)`
  border-left: 1px solid rgba(3, 24, 73, 0.1);
  padding: 24px;
  z-index: 2;
`;
