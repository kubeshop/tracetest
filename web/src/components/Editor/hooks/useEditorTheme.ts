import {tags as t} from '@lezer/highlight';
import {createTheme} from '@uiw/codemirror-themes';
import {useMemo} from 'react';
import {useTheme} from 'styled-components';

interface IProps {
  isEditable?: boolean;
}

const useEditorTheme = ({isEditable = true}: IProps = {}) => {
  const {
    color: {white, text, textLight},
  } = useTheme();

  return useMemo(
    () =>
      createTheme({
        theme: 'light',
        settings: {
          background: white,
          foreground: text,
          caret: text,
          selection: textLight,
          selectionMatch: white,
          lineHighlight: white,
          ...(isEditable
            ? {
                gutterBackground: white,
                gutterBorder: white,
              }
            : {
                gutterBackground: white,
                gutterBorder: white,
                gutterForeground: white,
              }),
        },
        styles: [
          {tag: t.tagName, color: '#994cc3'},
          {tag: t.string, color: '#4876d6'},
          {tag: t.operatorKeyword, color: '#994cc3'},
          {tag: t.keyword, color: text},
        ],
      }),
    [isEditable, text, textLight, white]
  );
};

export default useEditorTheme;
