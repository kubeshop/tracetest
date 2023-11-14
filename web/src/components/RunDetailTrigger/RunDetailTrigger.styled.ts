import styled from 'styled-components';

export const Body = styled.div`
  display: flex;
  height: calc(100vh - 171px);
  width: 100%;
`;

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
`;

export const Section = styled.div`
  background-color: ${({theme}) => theme.color.white};
  flex: 1;
  overflow-y: scroll;
  padding: 14px 24px;
`;

export const SectionLeft = styled(Section)`
  z-index: 1;
  flex-basis: 50%;
  max-width: 50vw;

  .ant-tabs-tab {
    padding-top: 0;
  }
`;

export const SectionRight = styled(Section)`
  border-left: 1px solid rgba(3, 24, 73, 0.1);
  z-index: 2;
  flex-basis: 50%;
  max-width: 50vw;
`;
