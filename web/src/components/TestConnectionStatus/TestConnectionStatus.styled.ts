import {CheckCircleFilled, InfoCircleFilled, LoadingOutlined, MinusCircleFilled} from '@ant-design/icons';
import styled from 'styled-components';

export const StatusContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

export const IconSuccess = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
`;

export const IconFail = styled(MinusCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const IconInfo = styled(InfoCircleFilled)`
  color: ${({theme}) => theme.color.textLight};
`;

export const LoadingIcon = styled(LoadingOutlined)`
  color: ${({theme}) => theme.color.primary};
`;
