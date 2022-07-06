import {CopyOutlined, PlusOutlined} from '@ant-design/icons';
import {Badge, Tooltip, Typography} from 'antd';
import styled from 'styled-components';

export const AttributeRow = styled.div`
  display: grid;
  grid-template-columns: 160px 1fr 60px;
  gap: 14px;
  align-items: start;
  padding: 8px;

  &:hover {
    background-color: ${({theme}) => theme.color.background};
  }
`;

export const AttributeValueRow = styled.div`
  display: flex;
  word-break: break-word;
`;

export const CopyIcon = styled(CopyOutlined)`
  color: ${({theme}) => theme.color.primary};
  cursor: pointer;
`;

export const AddAssertionIcon = styled(PlusOutlined)`
  color: ${({theme}) => theme.color.primary};
  cursor: pointer;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const Text = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.sm};
`;

export const IconContainer = styled.div`
  display: flex;
  gap: 14px;
  align-items: center;
  align-self: center;
`;

export const CustomTooltip = styled(Tooltip).attrs({
  placement: 'top',
  arrowPointAtCenter: true,
})``;

export const CustomBadge = styled(Badge)`
  border: ${({theme}) => `1px solid ${theme.color.textSecondary}`};
  border-radius: 9999px;
  cursor: pointer;
  line-height: 19px;
  margin-left: 8px;
  padding: 0 8px;
  white-space: nowrap;

  .ant-badge-status-text {
    color: ${({theme}) => theme.color.textSecondary};
    font-size: ${({theme}) => theme.size.sm};
    margin-left: 3px;
  }
`;
