import {CheckCircleOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled, {css} from 'styled-components';

export const CreateTestSteps = styled.div`
  height: inherit;
  position: relative;

  .ant-tabs {
    height: 100%;
  }

  .ant-tabs > .ant-tabs-nav .ant-tabs-nav-wrap {
    overflow: visible;
  }

  p {
    margin: 0;
  }

  .ant-tabs-nav {
    height: max-content;
  }

  .ant-tabs-nav-list {
    padding-left: 18px;
    display: flex;
    gap: 43px;
    flex-direction: column;
  }

  .ant-tabs-left > .ant-tabs-nav .ant-tabs-tab {
    padding: 8px 0;
  }

  > .ant-tabs-card .ant-tabs-content {
    height: 100%;
  }

  .ant-tabs-left > .ant-tabs-nav,
  .ant-tabs-right > .ant-tabs-nav,
  .ant-tabs-left > div > .ant-tabs-nav,
  .ant-tabs-right > div > .ant-tabs-nav {
    width: 160px;
  }

  > .ant-tabs-card .ant-tabs-content > .ant-tabs-tabpane {
    padding: 0;
    padding-left: 65px;
  }

  > .ant-tabs-card > .ant-tabs-nav::before {
    display: none;
  }

  > .ant-tabs-card .ant-tabs-tab {
    background: transparent;
    border: none;
  }

  > .ant-tabs-card .ant-tabs-tab-active {
    border: none;
  }

  .ant-tabs-card.ant-tabs-left > .ant-tabs-nav .ant-tabs-tab + .ant-tabs-tab {
    margin-top: 0;
  }

  .ant-tabs-card.ant-tabs-left > .ant-tabs-nav .ant-tabs-tab {
    border-radius: 0;
  }

  .ant-tabs-left > .ant-tabs-content-holder,
  .ant-tabs-left > div > .ant-tabs-content-holder {
    border: none;
  }
`;

export const StatusIcon = styled(CheckCircleOutlined)`
  color: ${({theme}) => theme.color.primary};
  font-size: ${({theme}) => theme.size.lg};
  margin-left: 5px;
`;

export const CreateStepsTab = styled.div`
  display: flex;
  align-items: center;
  word-break: normal;
  white-space: normal;
  text-align: left;
`;

export const CreateStepsTabTitle = styled(Typography.Text)<{$isActive: boolean}>`
  && {
    color: ${({theme}) => theme.color.primary};
    font-size: ${({theme}) => theme.size.md};
    font-weight: ${({$isActive}) => ($isActive ? '600' : '400')};
  }
`;

export const Step = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 426px;
  max-height: 426px;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 8px;
`;

export const FormContainer = styled.div`
  height: 100%;
`;

export const Title = styled(Typography.Title).attrs({
  level: 5,
})`
  && {
    margin-bottom: 24px;
  }
`;

export const LoadingSpinnerContainer = styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const ProgressLine = styled.div<{$stepCount: number}>`
  border: 1px solid ${({theme}) => theme.color.primary};
  height: calc(83px * ${({$stepCount}) => $stepCount - 1});
  position: absolute;
  top: 17px;
  left: 0;
`;

export const StepDot = styled.div<{$isActive: boolean}>`
  position: absolute;
  top: 17px;
  left: -20.5px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: ${({theme}) => theme.color.primary};

  ${({$isActive}) =>
    $isActive &&
    css`
      top: 15px;
      width: 11px;
      height: 11px;
      left: -22.5px;
    `}
`;
