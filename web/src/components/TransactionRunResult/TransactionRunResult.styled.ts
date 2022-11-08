import {CheckCircleFilled, CloseCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 2px;
  background: ${({theme}) => theme.color.background};
  display: flex;
  gap: 16px;
  padding: 7px 16px;
  margin-bottom: 8px;
  height: 58px;
`;

export const Info = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  width: 100%;
  height: 100%;
`;

export const Stack = styled.div`
  display: flex;
  justify-content: space-between;
  flex-direction: column;
`;

export const Text = styled(Typography.Text)<{opacity?: number}>`
  && {
    font-size: ${({theme}) => theme.size.sm};
    color: ${({opacity}) => `rgba(3, 24, 73, ${opacity || 1})`};
    margin: 0 !important;
  }
`;

export const Title = styled(Typography.Title)`
  font-size: ${({theme}) => theme.size.lg};
  margin-bottom: 24px;
`;

export const IconSuccess = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
`;

export const IconFail = styled(CloseCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const TagContainer = styled.div`
  display: flex;
  white-space: nowrap;
  overflow: auto;
`;
