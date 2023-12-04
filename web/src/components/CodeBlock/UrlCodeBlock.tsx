import * as S from './CodeBlock.styled';
import * as U from './UrlCodeBlock.styled';
import {IProps} from './CodeBlock';
import useCopy from '../../hooks/useCopy';

const UrlCodeBlock = ({...props}: IProps) => {
  const copy = useCopy();

  return (
    <S.FrameContainer>
      <U.CopyButton ghost type="primary" onClick={() => copy(props.value)}>
        Copy
      </U.CopyButton>
      <U.CodeBlock {...props} />
    </S.FrameContainer>
  );
};

export default UrlCodeBlock;
