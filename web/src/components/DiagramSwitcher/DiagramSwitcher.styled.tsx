import {BarsOutlined, ClusterOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const DiagramSwitcher = styled.div`
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 24px;
`;

export const Switch = styled.div`
  background: #fff;
  border: 1px solid rgba(3, 24, 73, 0.1);
  border-radius: 2px;
  display: flex;
  gap: 13px;
  padding: 7px;
`;

export const DAGIcon = styled(ClusterOutlined).attrs({
  style: {
    fontSize: '18px',
  },
})<{$isSelected?: boolean}>`
  cursor: pointer;
  color: ${({$isSelected = false}) => ($isSelected ? '#61175E' : '#9AA3AB')};
`;

export const TimelineIcon = styled(BarsOutlined).attrs<{isSelected?: boolean}>({
  style: {
    fontSize: '18px',
  },
})<{$isSelected?: boolean}>`
  cursor: pointer;
  color: ${({$isSelected = false}) => ($isSelected ? '#61175E' : '#9AA3AB')};
`;
