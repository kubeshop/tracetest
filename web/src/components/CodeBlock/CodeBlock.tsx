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

const getLang = (mimeType: string): LanguageName | undefined => langNames.find(lang => mimeType.includes(`/${lang}`));

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
  const lang = useMemo(() => language || getLang(mimeType), [language, mimeType]);

  return (
    <S.CodeContainer data-cy="code-block" $maxHeight={maxHeight} $minHeight={minHeight}>
      <SyntaxHighlighter language={lang} style={arduinoLight} wrapLongLines>
        {formatValue(value, lang)}
      </SyntaxHighlighter>
    </S.CodeContainer>
  );
};

export default CodeBlock;
