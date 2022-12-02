import {Badge} from 'antd';
import styled, {css} from 'styled-components';

export const Container = styled.div`
  display: flex;
  gap: 2px;
`;

export const CustomBadge = styled(Badge)<{$styleType: 'node' | 'summary' | 'default'}>`
  border: ${({theme}) => `1px solid ${theme.color.textSecondary}`};
  border-radius: 9999px;
  cursor: pointer;
  line-height: 19px;
  padding: 0 8px;
  white-space: nowrap;

  .ant-badge-status-text {
    color: ${({theme}) => theme.color.textSecondary};
    font-size: ${({theme}) => theme.size.sm};
    margin-left: 3px;
  }

  ${({$styleType}) =>
    $styleType === 'node' &&
    css`
      border: none;
      font-size: ${({theme}) => theme.size.xs};
      line-height: 10px;
      padding: 0 2px;
      vertical-align: bottom;
      margin: 0;

      .ant-badge-status-dot {
        height: 8px;
        top: 0;
        width: 8px;
      }

      .ant-badge-status-text {
        font-size: ${({theme}) => theme.size.xs};
        margin-left: 2px;
      }
    `}

  ${({$styleType}) =>
    $styleType === 'summary' &&
    css`
      border: none;
      padding: 0 2px;
      vertical-align: bottom;

      .ant-badge-status-dot {
        height: 10px;
        width: 10px;
      }
    `}
`;
