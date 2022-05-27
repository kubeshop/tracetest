import {CopyOutlined, PlusOutlined} from '@ant-design/icons';
import {Badge, Tooltip, Typography} from 'antd';
import styled from 'styled-components';

export const AttributeRow = styled.div`
  display: grid;
  grid-template-columns: 190px 1fr 60px;
  gap: 14px;
  align-items: start;
  padding: 8px;

  .ant-tooltip-inner {
    color: yellow;
  }

  &:hover {
    background-color: #fbfbff;
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
  align-self: center;
`;

export const CustomTooltip = styled(Tooltip).attrs({
  color: '#FBFBFF',
  placement: 'top',
  arrowPointAtCenter: true,
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
