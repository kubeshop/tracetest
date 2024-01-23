import {Typography} from 'antd';
import styled, {css} from 'styled-components';

export const Container = styled.div`
  margin: 8px 0;

  .ant-tabs-card > .ant-tabs-nav .ant-tabs-tab,
  .ant-tabs-card > div > .ant-tabs-nav .ant-tabs-tab {
    border: none;
    padding: 0;
  }

  .ant-tabs-tab .anticon {
    margin-right: 0;
  }

  .ant-tabs-left > .ant-tabs-content-holder > .ant-tabs-content > .ant-tabs-tabpane {
    padding: 16px;
  }

  .ant-tabs-left > .ant-tabs-content-holder,
  .ant-tabs-left > div > .ant-tabs-content-holder {
    margin-left: 0;
  }
`;

export const StepTabContainer = styled.div<{$isActive: boolean; $isDisabled: boolean}>`
  align-items: center;
  background-color: ${({theme, $isActive}) => ($isActive ? theme.color.backgroundInteractive : theme.color.white)};
  display: flex;
  gap: 8px;
  padding: 16px;
  text-align: left;
  min-width: 360px;

  ${({$isDisabled}) =>
    $isDisabled &&
    css`
      > div:first-child {
        border: 2px solid ${({theme}) => theme.color.textLight};
      }
      div,
      span {
        color: ${({theme}) => theme.color.textLight};
      }
    `}
`;

export const StepTabNumber = styled.div`
  border: 2px solid ${({theme}) => theme.color.textSecondary};
  border-radius: 50%;
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.sm};
  font-weight: 600;
  height: 24px;
  text-align: center;
  width: 24px;
`;

export const StepTabCheck = styled.div`
  align-items: center;
  background-color: ${({theme}) => theme.color.primary};
  border-radius: 50%;
  color: ${({theme}) => theme.color.white};
  display: flex;
  font-size: ${({theme}) => theme.size.sm};
  height: 24px;
  justify-content: center;
  width: 24px;
`;

export const StepTabTitle = styled(Typography.Text)<{$isActive: boolean}>`
  && {
    font-size: ${({theme}) => theme.size.lg};
    font-weight: ${({$isActive}) => ($isActive ? '600' : '400')};
  }
`;
