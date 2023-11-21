import styled, {css} from 'styled-components';
import {Typography} from 'antd';
import {LeftOutlined} from '@ant-design/icons';
import Link from 'components/Link';

export const BackIcon = styled(LeftOutlined)`
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;

export const Container = styled.div`
  height: 100vh;

  .run-tabs.ant-tabs,
  .run-tabs.ant-tabs .ant-tabs-content {
    height: 100%;
  }
`;

export const TransactionLink = styled(Link)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
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

    .ant-tabs-nav-wrap {
      overflow: visible;
    }

    .ant-tabs-nav-more {
      display: none;
    }

    ::before {
      display: none;
    }
  }

  .ant-tabs-nav-list {
    border-radius: 2px;
  }

  .ant-tabs-tab {
    font-weight: 600;
    padding: 0;
    margin: 7px 0;
    border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
    border-right: none;

    &:nth-of-type(1) {
      border-top-left-radius: 2px;
      border-bottom-left-radius: 2px;
    }

    &:nth-last-child(2) {
      border-right: ${({theme}) => `1px solid ${theme.color.borderLight}`};
      border-top-right-radius: 2px;
      border-bottom-right-radius: 2px;
    }

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

export const Row = styled.div<{$height?: number}>`
  display: flex;
  height: ${({$height}) => `${$height}px` || 'auto'};
`;

export const Section = styled.div<{$justifyContent: string}>`
  align-items: center;
  display: flex;
  gap: 14px;
  justify-content: ${({$justifyContent}) => $justifyContent};
  width: calc((100vw / 2) - 200px);
`;

export const StateContainer = styled.div`
  align-items: center;
  display: flex;
  justify-self: flex-end;
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

export const TabLink = styled(Link)<{$isActive: boolean; $isDisabled?: boolean}>`
  && {
    color: ${({theme, $isActive}) => $isActive && theme.color.white};
    padding: 5px 16px;

    &:hover,
    &:visited,
    &:focused {
      color: ${({theme, $isActive}) => $isActive && theme.color.white};
    }

    ${({$isDisabled}) =>
      $isDisabled &&
      css`
        opacity: 0.5;
        pointer-events: none;
      `}
  }
`;
