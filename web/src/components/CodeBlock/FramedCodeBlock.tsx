import * as S from './CodeBlock.styled';
import CodeBlock, {IProps as ICodeBlockProps} from './CodeBlock';
import useCopy from '../../hooks/useCopy';

interface IProps extends ICodeBlockProps {
  title: string;
  actions?: React.ReactNode;
}

const FramedCodeBlock = ({title, actions, ...props}: IProps) => {
  const copy = useCopy();

  return (
    <S.FrameContainer $isFullHeight={props.isFullHeight}>
      <S.FrameHeader>
        <S.FrameTitle>{title}</S.FrameTitle>
        <S.ActionsContainer>
          {actions}
          <S.CopyButton ghost type="primary" onClick={() => copy(props.value)}>
            Copy
          </S.CopyButton>
        </S.ActionsContainer>
      </S.FrameHeader>
      <CodeBlock {...props} />
    </S.FrameContainer>
  );
};

export default FramedCodeBlock;
