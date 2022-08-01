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
    background: ${({theme}) => theme.color.background};
    border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
    font-size: ${({theme}) => theme.size.sm};

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
