import * as S from './CodeBlock.styled';
import CodeBlock, {IProps as ICodeBlockProps} from './CodeBlock';
import useCopy from '../../hooks/useCopy';

interface IProps extends ICodeBlockProps {
  title: string;
}

const FramedCodeBlock = ({title, ...props}: IProps) => {
  const copy = useCopy();

  return (
    <S.FrameContainer>
      <S.FrameHeader>
        <S.FrameTitle>{title}</S.FrameTitle>
        <S.CopyButton ghost type="primary" onClick={() => copy(props.value)}>
          Copy
        </S.CopyButton>
      </S.FrameHeader>
      <CodeBlock {...props} />
    </S.FrameContainer>
  );
};

export default FramedCodeBlock;
