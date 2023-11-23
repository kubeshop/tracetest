import SyntaxHighlighter from 'react-syntax-highlighter';
import {arduinoLight} from 'react-syntax-highlighter/dist/cjs/styles/hljs';
import {LanguageName, langNames} from '@uiw/codemirror-extensions-langs';
import {useMemo} from 'react';
import * as S from './CodeBlock.styled';

export interface IProps {
  value: string;
  language?: string;
  mimeType?: string;
  maxHeight?: string;
  minHeight?: string;
}

const getLanguage = (mimeType: string): LanguageName | undefined => {
  const language = langNames.find(lang => mimeType.includes(`/${lang}`));
  // SyntaxHighlighter does not support html, so we need to use xml instead
  return language === 'html' ? 'xml' : language;
};

const formatValue = (value: string, lang?: string): string => {
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

const CodeBlock = ({value, language, mimeType = '', maxHeight = '', minHeight = ''}: IProps) => {
  const lang = useMemo(() => language || getLanguage(mimeType), [language, mimeType]);

  // SyntaxHighlighter has a performance problem, so we need to memoize it
  // See https://github.com/react-syntax-highlighter/react-syntax-highlighter/issues/302
  const memoizedHighlighter = useMemo(
    () => (
      <SyntaxHighlighter language={lang} style={arduinoLight} wrapLongLines>
        {formatValue(value, lang)}
      </SyntaxHighlighter>
    ),
    [lang, value]
  );

  return (
    <S.CodeContainer data-cy="code-block" $maxHeight={maxHeight} $minHeight={minHeight}>
      {memoizedHighlighter}
    </S.CodeContainer>
  );
};

export default CodeBlock;
