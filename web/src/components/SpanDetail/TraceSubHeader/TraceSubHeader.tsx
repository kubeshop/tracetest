import AnalyzerErrorsPopover from 'components/AnalyzerErrorsPopover';
import {IPropsSubHeader} from '../SpanDetail';
import * as S from './TraceSubHeader.styled';

const TraceSubHeader = ({analyzerErrors}: IPropsSubHeader) => {
  return analyzerErrors ? (
    <S.Container>
      <AnalyzerErrorsPopover errors={analyzerErrors} />
    </S.Container>
  ) : null;
};

export default TraceSubHeader;
