import {CopyOutlined, PlusOutlined} from '@ant-design/icons';
import {Tooltip, Typography} from 'antd';
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
