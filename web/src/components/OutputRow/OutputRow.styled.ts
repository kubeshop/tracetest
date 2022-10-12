import {MoreOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: ${({theme}) => theme.color.background};
  border: 1px solid ${({theme}) => theme.color.borderLight};
  padding: 7px 16px;
  transition: background-color 0.2s ease;
`;

export const Entry = styled.div``;

export const OutputDetails = styled.div`
  display: grid;
  grid-template-columns: 85px 1fr 1fr 1fr;
  align-items: center;
`;

export const Value = styled(Typography.Text)`
  display: flex;
  word-break: break-word;
  font-size: ${({theme}) => theme.size.sm};
  font-weight: 600;
`;

export const Key = styled(Typography.Text)`
  display: flex;
  word-break: break-word;
  color: ${({theme}) => theme.color.textLight};
  font-size: ${({theme}) => theme.size.sm};
  font-weight: 400;
`;

export const ActionButton = styled(MoreOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;
