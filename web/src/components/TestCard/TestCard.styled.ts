import {MoreOutlined} from '@ant-design/icons';
import {Button, Typography} from 'antd';
import styled from 'styled-components';

import emptyStateIcon from 'assets/SpanAssertionsEmptyState.svg';

export const TestCard = styled.div<{$isCollapsed: boolean}>`
  box-shadow: 0 4px 8px rgba(153, 155, 168, 0.1);
  background: ${({theme}) => theme.color.white};
  border-left: ${({$isCollapsed, theme}) => $isCollapsed && `2px solid ${theme.color.primary}`};
  border-radius: 2px;
`;

export const InfoContainer = styled.div`
  cursor: pointer;
  display: grid;
  align-items: center;
  grid-template-columns: 20px 1fr 60px 2fr 220px 100px 20px;
  gap: 24px;
  padding: 16px 24px;
`;

export const ResultListContainer = styled.div`
  margin: 0 68px 54px 54px;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const ButtonContainer = styled.div`
  display: flex;
  justify-content: flex-end;
`;

export const NameText = styled(Typography.Title).attrs({ellipsis: true, level: 3})`
  && {
    margin: 0;
  }
`;

export const Text = styled(Typography.Text)``;

export const ActionButton = styled(MoreOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;

export const TestDetails = styled.div`
  text-align: right;
  width: 100%;
  margin-top: 8px;
`;

export const TestDetailsLink = styled(Button).attrs({
  type: 'link',
})`
  color: ${({theme}) => theme.color.primary};
  font-weight: 600;
  padding: 0;
`;

export const EmptyStateIcon = styled.img.attrs({
  src: emptyStateIcon,
})``;

export const EmptyStateContainer = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: 14px;
  justify-content: center;
  margin: 10px 0 20px;
`;
