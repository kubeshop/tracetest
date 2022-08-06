import styled from 'styled-components';

export const AdvancedEditor = styled.div`
  //width: 100%;

  && {
    .cm-editor {
      border-radius: 2px;
      font-size: ${({theme}) => theme.size.md};
    }

    .cm-focused {
      outline: 1px solid ${({theme}) => theme.color.primary};
    }

    .cm-tooltip {
      border: none;
      background-color: ${({theme}) => theme.color.white};
      padding: 9px 48px 9px 24px;
      box-shadow: 0px 9px 28px 8px rgba(0, 0, 0, 0.05), 0px 6px 16px rgba(0, 0, 0, 0.08),
        0px 3px 6px -4px rgba(0, 0, 0, 0.12);
      border-radius: 2px;
    }

    .cm-scroller {
      scrollbar-width: none;
      -ms-overflow-style: none;
    }

    .cm-scroller::-webkit-scrollbar {
      display: none;
      -webkit-appearance: none;
      width: 0;
      height: 0;
    }
  }
`;
