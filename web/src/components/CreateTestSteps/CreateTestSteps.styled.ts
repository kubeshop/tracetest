import {CheckCircleOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';
import {TStepStatus} from 'types/Plugins.types';

export const CreateTestSteps = styled.div`
  margin: 24px;
  height: inherit;
  border: 1px solid rgba(3, 24, 73, 0.1);

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
    padding: 16px;
    background: #fff;
  }

  > .ant-tabs-card > .ant-tabs-nav::before {
    display: none;
  }

  > .ant-tabs-card .ant-tabs-tab {
    background: transparent;
    border-color: transparent;
    border-bottom: 1px solid rgba(3, 24, 73, 0.1);
  }

  > .ant-tabs-card .ant-tabs-tab-active {
    background: #fff;
    border-color: #fff;
    border-bottom: 1px solid rgba(3, 24, 73, 0.1);
  }

  .ant-tabs-card.ant-tabs-left > .ant-tabs-nav .ant-tabs-tab + .ant-tabs-tab {
    margin-top: 0;
  }

  .ant-tabs-card.ant-tabs-left > .ant-tabs-nav .ant-tabs-tab {
    border-radius: 0;
  }

  #components-tabs-demo-card-top .code-box-demo {
    padding: 24px;
    overflow: hidden;
    background: #f5f5f5;
  }
`;

const colorToStatusMap: Record<TStepStatus, string> = {
  complete: '#66BB6A',
  selected: '#031849',
  pending: 'rgba(3, 24, 73, 0.2);',
};

export const StatusIcon = styled(CheckCircleOutlined)<{$status?: TStepStatus}>`
  font-size: 20px;
  color: ${({$status = 'pending'}) => colorToStatusMap[$status]};
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
  justify-content: flex-end;
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
