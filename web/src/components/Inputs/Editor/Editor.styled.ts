import styled, {css} from 'styled-components';

export const EditorContainer = styled.div`
  width: 100%;

  && {
    .cm-editor {
      border-radius: 2px;
      font-size: ${({theme}) => theme.size.md};
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
    }

    .cm-content .cm-activeLine.cm-line::selection,
    .cm-content .cm-activeLine {
      color: inherit;
      background: inherit;
    }

    .cm-content span::selection {
      color: inherit;
      background: inherit;
    }

    .cm-content span::-moz-selection {
      /* Code for Firefox */
      color: inherit;
      background: inherit;
    }
  }
`;

export const SelectorEditorContainer = styled(EditorContainer)<{$isEditable: boolean}>`
  && {
    .cm-editor {
      outline: ${({$isEditable, theme}) => $isEditable && `1px solid ${theme.color.primary}`};
    }

    .cm-scroller {
      overflow: inherit;
    }

    ${({$isEditable}) =>
      !$isEditable &&
      css`
        .cm-gutterElement {
          display: none;
        }
      `}
  }
`;

export const ExpressionEditorContainer = styled(EditorContainer)<{$isEditable: boolean}>`
  && {
    .cm-editor {
      border: 1px solid ${({theme}) => theme.color.border};
      transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);

      &.cm-focused {
        border-color: ${({theme}) => theme.color.primary};
        box-shadow: 0 0 0 2px rgb(97 23 94 / 20%);
        border-right-width: 1px;
        outline: 0;
      }

      .cm-content {
        padding: 5px 0;
      }
    }

    .cm-content {
      font-family: ${({theme}) => theme.font.family};
    }

    ${({$isEditable}) =>
      !$isEditable &&
      css`
        .cm-editor {
          border: none;
        }

        .cm-gutterElement {
          display: none;
        }

        .cm-gutters {
          border: none;
        }

        .cm-line {
          padding: 0;
        }
      `}
  }
`;

export const InterpolationEditorContainer = styled(EditorContainer)<{$showLineNumbers?: boolean}>`
  && {
    ${({$showLineNumbers = false}) =>
      !$showLineNumbers &&
      css`
        .cm-gutterElement {
          display: none;
        }

        .cm-editor {
          height: 32px;
        }
      `}

    .cm-editor {
      border: 1px solid ${({theme}) => theme.color.border};
      transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);

      &.cm-focused {
        border-color: ${({theme}) => theme.color.primary};
        box-shadow: 0 0 0 2px rgb(97 23 94 / 20%);
        border-right-width: 1px;
        outline: 0;
      }
    }

    .cm-content {
      font-family: ${({theme}) => theme.font.family};
    }
  }
`;
