import {InfoCircleFilled, MoreOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div<{$isDeleted: boolean; $isSelected: boolean}>`
  background: ${({theme}) => theme.color.white};
  border: ${({$isSelected, theme}) =>
    $isSelected ? `1px solid ${theme.color.interactive}` : `1px solid ${theme.color.borderLight}`};
  display: flex;
  flex-direction: column;
  padding: 7px 16px;
  transition: background-color 0.2s ease;
  cursor: pointer;

  > div:nth-child(2) {
    opacity: ${({$isDeleted}) => ($isDeleted ? 0.5 : 1)};
  }
`;

export const Row = styled.div<{$justifyContent?: string}>`
  align-items: center;
  display: flex;
  justify-content: ${({$justifyContent}) => $justifyContent && $justifyContent};
`;

export const Entry = styled.div``;

export const OutputDetails = styled.div`
  align-items: flex-start;
  column-gap: 8px;
  display: grid;
  flex: 1;
  grid-template-columns: 1fr 2fr 1fr;
  margin-bottom: 8px;
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

export const ActionsContainer = styled.div`
  align-items: center;
  display: flex;
  justify-content: flex-end;
`;

export const IconWarning = styled(InfoCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;
