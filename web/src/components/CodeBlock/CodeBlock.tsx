import ReactCodeMirror from '@uiw/react-codemirror';
import {loadLanguage, LanguageName} from '@uiw/codemirror-extensions-langs';
import {useMemo} from 'react';
import {isHTML, isJson} from 'utils/Common';
import * as S from './CodeBlock.styled';

interface IProps {
  value: string;
}

const getInitialLang = (value: string): LanguageName | undefined => {
  if (isJson(value)) return 'json';
  if (isHTML(value)) return 'html';

  return undefined;
};

const formatValue = (value: string, lang: LanguageName | undefined): string => {
  switch (lang) {
    case 'json':
      return JSON.stringify(JSON.parse(value), null, 2);

    default:
      return value;
  }
};

const CodeBlock = ({value}: IProps) => {
  const lang = useMemo(() => getInitialLang(value), [value]);
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
