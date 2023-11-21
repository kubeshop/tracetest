import styled from 'styled-components';

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
`;

export const Header = styled.div`
  align-items: flex-start;
  background-color: ${({theme}) => theme.color.white};
  display: flex;
  gap: 20px;
  padding: 18px 24px 14px;
`;

export const HeaderLeft = styled.div`
  flex: 1;
  overflow: hidden;

  .ant-form-item-with-help .ant-form-item-explain {
    display: none;
    min-height: 0;
    height: 0;
    transition: none;
  }
`;

export const HeaderRight = styled.div`
  width: 110px;
`;

export const Body = styled.div`
  display: flex;
  height: calc(100vh - 172px);
  width: 100%;
`;

export const Section = styled.div`
  background-color: ${({theme}) => theme.color.white};
  flex: 1;
  overflow-y: scroll;
  padding: 14px 24px;
`;

export const SectionLeft = styled(Section)`
  .ant-tabs-tab {
    padding-top: 0;
  }
`;

export const SectionRight = styled(Section)`
  border-left: 1px solid rgba(3, 24, 73, 0.1);
`;

export const EmptyContainer = styled.div`
  margin-top: 132px;

  .ant-empty-description {
    color: ${({theme}) => theme.color.textSecondary};
  }
`;
