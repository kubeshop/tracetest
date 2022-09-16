import styled from 'styled-components';

export const ExpectedInputContainer = styled.div`
  width: 0;
  flex-basis: 50%;
  padding-left: 8px;

  && {
    .cm-editor {
      overflow: hidden;
      display: flex;
      border-radius: 2px;
      font-size: ${({theme}) => theme.size.md};
      outline: 1px solid grey;
      height: 32px;
      font-family: SFPro, serif;
    }
    .cm-content {
      display: flex;
      align-items: center;
    }
    .cm-scroller {
      overflow: hidden;
    }
    .cm-line {
      padding: 0;
      span {
        font-family: SFPro, serif;
      }
    }
  }
`;
