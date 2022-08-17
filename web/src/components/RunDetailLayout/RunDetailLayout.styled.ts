import styled from 'styled-components';
import {Typography} from 'antd';
import {LeftOutlined} from '@ant-design/icons';

export const BackIcon = styled(LeftOutlined)`
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;

export const Container = styled.div`
  height: 100%;

  .ant-tabs,
  .ant-tabs-content {
    height: 100%;
  }
`;

export const ContainerHeader = styled.div`
  background-color: ${({theme}) => theme.color.white};
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  padding: 6px 24px;
  width: 100%;

  && .ant-tabs-nav {
    margin-bottom: 0;

    .ant-tabs-ink-bar {
      display: none;
    }

    ::before {
      display: none;
    }
  }

  .ant-tabs-tab {
    border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
    font-weight: 600;
    padding: 5px 16px;
    margin: 7px 0;

    &.ant-tabs-tab-active {
      background-color: ${({theme}) => theme.color.primary};

      .ant-tabs-tab-btn {
        color: ${({theme}) => theme.color.white};
      }
    }
  }
`;

export const InfoContainer = styled.div`
  flex: 1;
`;

export const Row = styled.div`
  display: flex;
`;

export const Section = styled.div<{$justifyContent: string}>`
  align-items: center;
  display: flex;
  gap: 14px;
  justify-content: ${({$justifyContent}) => $justifyContent};
  // width: 250px;
  width: calc((100vw / 2) - 150px);
`;

export const StateContainer = styled.div`
  align-items: center;
  display: flex;
  justify-self: flex-end;
  cursor: pointer;
`;

export const StateText = styled(Typography.Text)`
  && {
    margin-right: 8px;
    color: ${({theme}) => theme.color.textSecondary};
  }
`;

export const Text = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  && {
    font-size: ${({theme}) => theme.size.sm};
    margin: 0;
  }
`;

export const Title = styled(Typography.Title).attrs({ellipsis: true, level: 2})`
  && {
    margin: 0;
    // max-width: calc(100% - 24px);
    max-width: calc((100vw / 2) - 150px - 54px);
  }
`;
