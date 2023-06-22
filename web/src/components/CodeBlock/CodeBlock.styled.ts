import {Button, Typography} from 'antd';
import styled from 'styled-components';

export const CodeContainer = styled.div<{$maxHeight: string; $minHeight: string}>`
  position: relative;
  border: ${({theme}) => `1px solid ${theme.color.border}`};
  min-height: ${({$minHeight}) => $minHeight || '370px'};

  pre {
    margin: 0;
    padding: 13px 16px !important;
    min-height: inherit;
    max-height: ${({$maxHeight}) => $maxHeight || '340px'};
    background: ${({theme}) => theme.color.background} !important;

    &:hover {
      background: ${({theme}) => theme.color.backgroundInteractive} !important;
    }
  }
`;

export const FrameContainer = styled.div``;
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

export const CopyButton = styled(Button)`
  && {
    padding: 0 8px;
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
