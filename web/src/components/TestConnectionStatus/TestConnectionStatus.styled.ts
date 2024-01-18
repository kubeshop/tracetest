import {CheckCircleFilled, InfoCircleFilled, LoadingOutlined, MinusCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
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

export const StatusText = styled(Typography.Paragraph).attrs({
  type: 'secondary',
})`
  && {
    font-weight: 600;
    margin: 0;
  }
`;

export const SpanCountTest = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  && {
    font-size: ${({theme}) => theme.size.sm};
  }
`;
