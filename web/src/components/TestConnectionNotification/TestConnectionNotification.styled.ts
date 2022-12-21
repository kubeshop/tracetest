import {CheckCircleFilled, MinusCircleFilled} from '@ant-design/icons';
import styled from 'styled-components';

export const SuccessCheckIcon = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
  margin-top: 5px;
`;
export const FailedCheckIcon = styled(MinusCircleFilled)`
  color: ${({theme}) => theme.color.error};
  margin-top: 5px;
`;
export const Container = styled.div``;
export const StepContainer = styled.div`
  display: flex;
  gap: 6px;
  align-items: start;
`;
