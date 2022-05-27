import {CopyOutlined, PlusOutlined} from '@ant-design/icons';
import {Badge, Tooltip, Typography} from 'antd';
import styled from 'styled-components';

export const AttributeRow = styled.div`
  display: grid;
  grid-template-columns: 200px 1fr 50px;
  gap: 14px;
  align-items: start;

  .ant-tooltip-inner {
    color: yellow;
  }
`;

export const AttributeValueRow = styled.div`
  display: flex;
`;

export const CopyIcon = styled(CopyOutlined)`
  cursor: pointer;
  color: #61175e;
`;

export const AddAssertionIcon = styled(PlusOutlined)`
  cursor: pointer;
  color: #61175e;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const Text = styled(Typography.Text)`
  font-size: 12px;
`;

export const IconContainer = styled.div`
  display: flex;
  gap: 14px;
  align-items: center;
`;

export const CustomTooltip = styled(Tooltip).attrs({
  color: '#FBFBFF',
})``;

export const CustomBadge = styled(Badge)`
  border: 1px solid rgba(3, 24, 73, 0.5);
  border-radius: 9999px;
  cursor: pointer;
  line-height: 19px;
  margin-left: 8px;
  padding: 0 8px;
  white-space: nowrap;

  .ant-badge-status-text {
    color: rgba(3, 24, 73, 0.5);
    font-size: 12px;
    margin-left: 3px;
  }
`;
