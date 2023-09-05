import styled from 'styled-components';

export const BodyFieldContainer = styled.div<{$isDisplaying: boolean}>`
  width: 100%;
  display: ${({$isDisplaying}) => ($isDisplaying ? 'none' : 'unset')};

  && {
    .cm-editor {
      overflow: hidden;
      display: flex;
      border-radius: 2px;
      font-size: ${({theme}) => theme.size.md};
      outline: 1px solid grey;
      font-family: SFPro, Inter, serif;
    }

    .cm-line {
      padding: 0;

      span {
        font-family: SFPro, Inter, serif;
      }
    }
  }
`;
