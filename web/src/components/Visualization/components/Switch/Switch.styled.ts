import {BarsOutlined, ClusterOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const Container = styled.div`
  background: ${({theme}) => theme.color.white};
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 2px;
  display: flex;
  gap: 12px;
  padding: 7px;
`;

export const DAGIcon = styled(ClusterOutlined)<{$isSelected?: boolean}>`
  color: ${({$isSelected = false, theme}) => ($isSelected ? theme.color.primary : theme.color.textSecondary)};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.xl};
`;

export const TimelineIcon = styled(BarsOutlined)<{$isSelected?: boolean}>`
  color: ${({$isSelected = false, theme}) => ($isSelected ? theme.color.primary : theme.color.textSecondary)};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.xl};
`;
