import {BarsOutlined, ClusterOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const DiagramSwitcher = styled.div`
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 24px;
`;

export const Switch = styled.div`
  background: ${({theme}) => theme.color.white};
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 2px;
  display: flex;
  gap: 13px;
  padding: 7px;
`;

export const DAGIcon = styled(ClusterOutlined)<{$isSelected?: boolean}>`
  cursor: pointer;
  color: ${({$isSelected = false, theme}) => ($isSelected ? theme.color.primary : theme.color.textSecondary)};
  font-size: ${({theme}) => theme.size.xl};
`;

export const TimelineIcon = styled(BarsOutlined)<{$isSelected?: boolean}>`
  cursor: pointer;
  color: ${({$isSelected = false, theme}) => ($isSelected ? theme.color.primary : theme.color.textSecondary)};
  font-size: ${({theme}) => theme.size.xl};
`;
