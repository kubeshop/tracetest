import {capitalize} from 'lodash';
import {SemanticGroupNames, SemanticGroupNamesToText} from '../../constants/SemanticGroupNames.constants';
import SpanService from '../../services/Span.service';
import {ISpan} from '../../types/Span.types';
import Generic from './components/Generic';
import Http from './components/Http';
import SpanHeader from './SpanHeader';
import * as S from './SpanDetail.styled';
import {useAppSelector} from '../../redux/hooks';
import {IAssertionResultList} from '../../types/Assertion.types';
import AssertionSelectors from '../../selectors/Assertion.selectors';

export interface ISpanDetailProps {
  testId?: string;
  span?: ISpan;
  resultId?: string;
  assertionsResultList?: IAssertionResultList[];
}

const ComponentMap: Record<string, typeof Generic> = {
  [SemanticGroupNames.Http]: Http,
};

const getSpanTitle = (span: ISpan) => {
  const {primary, heading} = SpanService.getSpanNodeInfo(span);
  const spanTypeText = SemanticGroupNamesToText[span.type];

  return `${capitalize(heading) || spanTypeText} • ${primary} • ${span.name}`;
};

const SpanDetail: React.FC<ISpanDetailProps> = ({span, testId, resultId}) => {
  const Component = ComponentMap[span?.type || ''] || Generic;

  const title = (span && getSpanTitle(span)) || '';
  const assertionsResultList = useAppSelector(
    AssertionSelectors.selectAssertionResultListBySpan(testId, resultId, span?.spanId)
  );

  return (
    <S.SpanDetail>
      <SpanHeader title={title} />
      <Component span={span} assertionsResultList={assertionsResultList} testId={testId} resultId={resultId} />
    </S.SpanDetail>
  );
};

export default SpanDetail;
