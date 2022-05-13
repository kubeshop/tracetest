import {CopyOutlined, PlusOutlined} from '@ant-design/icons';
import {Tooltip, Typography} from 'antd';
import styled, {css} from 'styled-components';

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

export const ValueText = styled(Text)<{isCollapsed: boolean}>`
  cursor: pointer;

  ${({isCollapsed}) =>
    isCollapsed
      ? css`
          word-break: break-all;
        `
      : css`
          display: -webkit-box;
          -webkit-line-clamp: 3;
          -webkit-box-orient: vertical;
          overflow: hidden;
        `}

  pre {
    margin: 0;
    background: #fbfbff;
    border: 1px solid #c9cedb;
    ${({isCollapsed}) =>
      isCollapsed
        ? css`
            word-break: break-all;
          `
        : css`
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
          `}
  }
`;

export const IconContainer = styled.div`
  display: flex;
  gap: 14px;
  align-items: center;
`;

export const CustomTooltip = styled(Tooltip).attrs({
  color: '#FBFBFF',
})``;
