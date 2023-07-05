import {ExclamationCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  background-color: ${({theme}) => theme.color.white};
  display: flex;
  padding: 0 12px;
`;

export const ErrorIcon = styled(ExclamationCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const Text = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.sm};
`;
