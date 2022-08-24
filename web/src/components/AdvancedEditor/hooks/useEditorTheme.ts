import {tags as t} from '@lezer/highlight';
import {createTheme} from '@uiw/codemirror-themes';
import {useMemo} from 'react';
import {useTheme} from 'styled-components';

interface IProps {
  isEditable?: boolean;
}

const useEditorTheme = ({isEditable = true}: IProps = {}) => {
  const {
    color: {white, text},
  } = useTheme();

  return useMemo(
    () =>
      createTheme({
        theme: 'light',
        settings: {
          background: white,
          foreground: text,
          caret: text,
          selection: text,
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
          {tag: t.tagName, color: '#66BB6A'},
          {tag: t.string, color: '#F03950'},
          {tag: t.operatorKeyword, color: text},
          {tag: t.keyword, color: text},
        ],
      }),
    [isEditable, text, white]
  );
};

export default useEditorTheme;
