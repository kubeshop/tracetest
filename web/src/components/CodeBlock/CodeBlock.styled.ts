import styled from 'styled-components';

export const CodeContainer = styled.div<{$maxHeight: string; $minHeight: string}>`
  position: relative;
  border: ${({theme}) => `1px solid ${theme.color.border}`};
  min-height: ${({$minHeight}) => $minHeight || '370px'};
  cursor: pointer;

  pre {
    margin: 0;
    min-height: inherit;
    max-height: ${({$maxHeight}) => $maxHeight || '340px'};
    background: ${({theme}) => theme.color.background} !important;

    &:hover {
      background: ${({theme}) => theme.color.backgroundInteractive} !important;
    }
  }
`;
