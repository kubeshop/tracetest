import ReactCodeMirror from '@uiw/react-codemirror';
import {loadLanguage, LanguageName, langNames} from '@uiw/codemirror-extensions-langs';
import {useMemo} from 'react';
import * as S from './CodeBlock.styled';

interface IProps {
  value: string;
  mimeType: string;
}

const getInitialLang = (mimeType: string): LanguageName | undefined =>
  langNames.find(lang => mimeType.includes(`/${lang}`));

const formatValue = (value: string, lang: LanguageName | undefined): string => {
  switch (lang) {
    case 'json':
      return JSON.stringify(JSON.parse(value), null, 2);

    default:
      return value;
  }
};

const CodeBlock = ({value, mimeType}: IProps) => {
  const lang = useMemo(() => getInitialLang(mimeType), [mimeType]);
  const extensionList = useMemo(() => (lang ? [loadLanguage(lang)!] : []), [lang]);

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
