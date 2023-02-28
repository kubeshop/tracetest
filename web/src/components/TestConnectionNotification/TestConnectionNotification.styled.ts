import {CheckCircleFilled, MinusCircleFilled, WarningFilled} from '@ant-design/icons';
import styled from 'styled-components';
import {Typography} from 'antd';

export const SuccessCheckIcon = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
  margin-top: 3px;
`;

export const FailedCheckIcon = styled(MinusCircleFilled)`
  color: ${({theme}) => theme.color.error};
  margin-top: 3px;
`;

export const WarningCheckIcon = styled(WarningFilled)`
  color: ${({theme}) => theme.color.warningYellow};
  margin-top: 3px;
`;

export const Container = styled.div``;

export const StepContainer = styled.div`
  align-items: start;
  display: flex;
  gap: 6px;
  margin-bottom: 8px;
`;

export const Title = styled(Typography.Title)`
  && {
    margin: 0;
  }
`;
