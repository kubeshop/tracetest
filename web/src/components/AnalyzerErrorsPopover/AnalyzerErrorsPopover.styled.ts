import {ExclamationCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  display: flex;
  gap: 4px;
`;

export const ContentContainer = styled.div`
  overflow-wrap: break-word;
  width: 270px;
`;

export const ErrorIcon = styled(ExclamationCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const List = styled.ul`
  padding-inline-start: 20px;
`;

export const RuleContainer = styled.div`
  margin-bottom: 8px;
`;

export const Text = styled(Typography.Text)``;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 0;
  }
`;
