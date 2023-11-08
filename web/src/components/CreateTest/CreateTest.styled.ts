import styled from 'styled-components';

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
`;

export const Header = styled.div`
  align-items: flex-start;
  display: flex;
  gap: 20px;
  padding: 18px 24px 14px;
`;

export const HeaderLeft = styled.div`
  flex: 1;
  overflow: hidden;
`;

export const HeaderRight = styled.div`
  width: 110px;
`;

export const Body = styled.div`
  display: flex;
  height: calc(100vh - 211px);
  width: 100%;
`;

export const Section = styled.div`
  flex: 1;
  padding: 14px 24px;
`;

export const SectionLeft = styled(Section)`
  background-color: ${({theme}) => theme.color.white};
  overflow-y: scroll;
`;

export const SectionRight = styled(Section)`
  background-color: ${({theme}) => theme.color.white};
  border-left: 1px solid rgba(3, 24, 73, 0.1);
  overflow-y: scroll;
`;

export const EmptyContainer = styled.div`
  margin-top: 132px;

  .ant-empty-description {
    color: ${({theme}) => theme.color.textSecondary};
  }
`;
