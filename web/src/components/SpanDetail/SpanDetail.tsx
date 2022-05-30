import {capitalize} from 'lodash';
import {useCallback} from 'react';
import {SemanticGroupNamesToText} from '../../constants/SemanticGroupNames.constants';
import SpanService from '../../services/Span.service';
import {TSpan, TSpanFlatAttribute} from '../../types/Span.types';
import SpanHeader from './SpanHeader';
import * as S from './SpanDetail.styled';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {CompareOperator} from '../../constants/Operator.constants';
import OperatorService from '../../services/Operator.service';
import SpanDetailTabs from './SpanDetailTabs';

export interface TSpanDetailProps {
  testId?: string;
  span?: TSpan;
  resultId?: string;
}

export interface TSpanDetailsComponentProps {
  span?: TSpan;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
}

const getSpanTitle = (span: TSpan) => {
  const {primary, heading} = SpanService.getSpanNodeInfo(span);
  const spanTypeText = SemanticGroupNamesToText[span.type];

  return `${capitalize(heading) || spanTypeText} • ${primary} • ${span.name}`;
};

const SpanDetail: React.FC<TSpanDetailProps> = ({span}) => {
  const {open} = useAssertionForm();

  const title = (span && getSpanTitle(span)) || '';

  const onCreateAssertion = useCallback(
    ({value, key}: TSpanFlatAttribute) => {
      const {selectorList, pseudoSelector} = SpanService.getSelectorInformation(span!);

      open({
        isEditing: false,
        defaultValues: {
          pseudoSelector,
          assertionList: [
            {
              comparator: OperatorService.getOperatorSymbol(CompareOperator.EQUALS),
              expected: value,
              attribute: key,
            },
          ],
          selectorList,
        },
      });
    },
    [open, span]
  );

  return (
    <S.SpanDetail>
      <SpanHeader title={title} />
      <SpanDetailTabs onCreateAssertion={onCreateAssertion} span={span} />
    </S.SpanDetail>
  );
};

export default SpanDetail;
