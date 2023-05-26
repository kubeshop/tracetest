import ReactCodeMirror from '@uiw/react-codemirror';
import {loadLanguage, LanguageName, langNames} from '@uiw/codemirror-extensions-langs';
import {useMemo} from 'react';
import * as S from './CodeBlock.styled';
import useEditorTheme from '../Editor/hooks/useEditorTheme';

interface IProps {
  value: string;
  mimeType: string;
}

const getInitialLang = (mimeType: string): LanguageName | undefined =>
  langNames.find(lang => mimeType.includes(`/${lang}`));

const formatValue = (value: string, lang: LanguageName | undefined): string => {
  switch (lang) {
    case 'json':
      try {
        return JSON.stringify(JSON.parse(value), null, 2);
      } catch (error) {
        return '';
      }

    default:
      return value;
  }
};

const CodeBlock = ({value, mimeType}: IProps) => {
  const lang = useMemo(() => getInitialLang(mimeType), [mimeType]);
  const theme = useEditorTheme();
  const extensionList = useMemo(() => (lang ? [theme, loadLanguage(lang)!] : []), [lang, theme]);

  return (
    <S.Container>
      <ReactCodeMirror
        basicSetup={{lineNumbers: false, foldGutter: false}}
        id="code-block"
        data-cy="code-block"
        value={formatValue(value, lang)}
        readOnly
        extensions={extensionList}
        spellCheck={false}
        autoFocus
      />
    </S.Container>
  );
};

export default CodeBlock;
