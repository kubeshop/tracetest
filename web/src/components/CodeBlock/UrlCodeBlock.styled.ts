import styled from 'styled-components';
import DefaultCodeBlock from './CodeBlock';
import * as S from './CodeBlock.styled';

export const CodeBlock = styled(DefaultCodeBlock)`
  && {
    pre {
      padding: 5px 10px !important;
      border-radius: 4px;
    }
  }
`;

export const CopyButton = styled(S.CopyButton)`
  && {
    position: absolute;
    top: 5px;
    right: 15px;
    z-index: 1;
  }
`;
