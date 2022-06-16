import {Typography} from 'antd';
import styled, {css} from 'styled-components';

export const ValueJson = styled(Typography.Text)<{$isCollapsed: boolean}>`
  cursor: pointer;

  ${({$isCollapsed}) =>
    $isCollapsed
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
    font-size: 12px;

    ${({$isCollapsed}) =>
      $isCollapsed
        ? css`
            word-break: break-all;
          `
        : css`
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
            word-break: break-all;
          `}
  }
`;

export const ValueText = styled(Typography.Text)``;
