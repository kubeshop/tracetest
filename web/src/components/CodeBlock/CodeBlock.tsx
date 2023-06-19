import SyntaxHighlighter from 'react-syntax-highlighter';
import {arduinoLight} from 'react-syntax-highlighter/dist/cjs/styles/hljs';
import {LanguageName, langNames} from '@uiw/codemirror-extensions-langs';
import useCopy from 'hooks/useCopy';
import {useMemo} from 'react';
import * as S from './CodeBlock.styled';

interface IProps {
  value: string;
  language?: string;
  mimeType?: string;
  maxHeight?: string;
  minHeight?: string;
}

const getLang = (mimeType: string): LanguageName | undefined => langNames.find(lang => mimeType.includes(`/${lang}`));

const CodeBlock = ({value, language, mimeType = '', maxHeight = '', minHeight = ''}: IProps) => {
  const copy = useCopy();
  const lang = useMemo(() => language || getLang(mimeType), [language, mimeType]);

  return (
    <S.CodeContainer data-cy="code-block" $maxHeight={maxHeight} $minHeight={minHeight} onClick={() => copy(value)}>
      <SyntaxHighlighter language={lang} style={arduinoLight} wrapLongLines>
        {value}
      </SyntaxHighlighter>
    </S.CodeContainer>
  );
};

export default CodeBlock;
