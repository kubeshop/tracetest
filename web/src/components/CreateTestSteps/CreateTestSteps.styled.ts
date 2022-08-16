import {CheckCircleOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';
import {TStepStatus} from 'types/Plugins.types';

export const CreateTestSteps = styled.div`
  margin: 24px;
  height: inherit;
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};

  .ant-tabs {
    height: 100%;
  }

  p {
    margin: 0;
  }

  > .ant-tabs-card .ant-tabs-content {
    height: 100%;
  }

  > .ant-tabs-card .ant-tabs-content > .ant-tabs-tabpane {
    background: ${({theme}) => theme.color.white};
    padding: 16px;
  }

  > .ant-tabs-card > .ant-tabs-nav::before {
    display: none;
  }

  > .ant-tabs-card .ant-tabs-tab {
    background: transparent;
    border-color: transparent;
    border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  }

  > .ant-tabs-card .ant-tabs-tab-active {
    background: ${({theme}) => theme.color.white};
    border-color: ${({theme}) => theme.color.white};
    border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  }

  .ant-tabs-card.ant-tabs-left > .ant-tabs-nav .ant-tabs-tab + .ant-tabs-tab {
    margin-top: 0;
  }

  .ant-tabs-card.ant-tabs-left > .ant-tabs-nav .ant-tabs-tab {
    border-radius: 0;
  }
`;

export const StatusIcon = styled(CheckCircleOutlined)<{$status?: TStepStatus}>`
  color: ${({$status, theme}) => ($status === 'complete' ? theme.color.success : theme.color.text)};
  font-size: 20px;
  opacity: ${({$status = 'pending'}) => $status === 'pending' && '0.2'};
`;

export const CreateTestStepsTab = styled.div`
  padding: 16px 2px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-width: 310px;
`;

export const Step = styled.div`
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
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
