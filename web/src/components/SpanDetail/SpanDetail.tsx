import {capitalize} from 'lodash';
import {useCallback} from 'react';
import {SemanticGroupNames, SemanticGroupNamesToText} from '../../constants/SemanticGroupNames.constants';
import {useCreateAssertionModal} from '../CreateAssertionModal/CreateAssertionModalProvider';
import SpanService from '../../services/Span.service';
import {ISpan, ISpanFlatAttribute} from '../../types/Span.types';
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
}

export interface ISpanDetailsComponentProps extends ISpanDetailProps {
  assertionsResultList: IAssertionResultList[];
  onCreateAssertion(attribute: ISpanFlatAttribute): void;
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
  const {open} = useCreateAssertionModal();
  const Component = ComponentMap[span?.type || ''] || Generic;

  const title = (span && getSpanTitle(span)) || '';
  const assertionsResultList = useAppSelector(
    AssertionSelectors.selectAssertionResultListBySpan(testId, resultId, span?.spanId)
  );

  const onCreateAssertion = useCallback(
    (attribute: ISpanFlatAttribute) => {
      open({
        span,
        resultId: resultId!,
        testId: testId!,
        defaultAttributeList: [attribute],
      });
    },
    [open, resultId, span, testId]
  );

  return (
    <S.SpanDetail>
      <SpanHeader title={title} />
      <Component
        span={span}
        assertionsResultList={assertionsResultList}
        testId={testId}
        resultId={resultId}
        onCreateAssertion={onCreateAssertion}
      />
    </S.SpanDetail>
  );
};

export default SpanDetail;
