import {Button, Typography} from 'antd';
import styled, {css} from 'styled-components';

export const CodeContainer = styled.div<{$maxHeight: string; $minHeight: string; $isFullHeight?: boolean}>`
  position: relative;
  min-height: ${({$minHeight}) => $minHeight || '370px'};

  pre {
    margin: 0;
    padding: 13px 16px !important;
    min-height: inherit;
    border: ${({theme}) => `1px solid ${theme.color.border}`};
    max-height: ${({$maxHeight}) => $maxHeight || '340px'};
    background: ${({theme}) => theme.color.background} !important;

    &:hover {
      background: ${({theme}) => theme.color.backgroundInteractive} !important;
    }
  }

  ${({$isFullHeight}) =>
    $isFullHeight &&
    css`
      height: calc(100% - 49px);

      pre {
        max-height: 100%;
      }
    `}
`;

export const FrameContainer = styled.div<{$isFullHeight?: boolean}>`
  position: relative;

  ${({$isFullHeight}) =>
    $isFullHeight &&
    css`
      height: calc(100% - 72px);
    `}
`;

export const FrameHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background: rgba(97, 23, 94, 0.1);
  border-radius: 2px 2px 0 0;
  border: 1px solid ${({theme}) => theme.color.border};
  border-bottom: none;
`;
export const FrameTitle = styled(Typography.Paragraph)`
  && {
    margin: 0;
  }
`;

export const CopyButton = styled(Button).attrs({
  size: 'small',
})`
  && {
    padding: 0 8px;
    width: 80px;
    background: ${({theme}) => theme.color.white};
    font-weight: 600;

    &:hover,
    &:focus,
    &:active {
      background: ${({theme}) => theme.color.white};
    }
  }
`;

export const ActionsContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;
