import {IPropsSubHeader} from '../SpanDetail';
import * as S from './TraceSubHeader.styled';

const TraceSubHeader = ({analyzerErrors}: IPropsSubHeader) => {
  return analyzerErrors ? (
    <S.Container>
      <S.ErrorIcon />
      <S.Text type="secondary">Analyzer errors</S.Text>
    </S.Container>
  ) : null;
};

export default TraceSubHeader;
