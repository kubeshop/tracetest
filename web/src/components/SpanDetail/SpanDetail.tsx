import {capitalize} from 'lodash';
import {useCallback} from 'react';
import {SemanticGroupNames, SemanticGroupNamesToText} from '../../constants/SemanticGroupNames.constants';
import SpanService from '../../services/Span.service';
import {TSpan, TSpanFlatAttribute} from '../../types/Span.types';
import Generic from './components/Generic';
import Http from './components/Http';
import SpanHeader from './SpanHeader';
import * as S from './SpanDetail.styled';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {CompareOperator} from '../../constants/Operator.constants';

export interface TSpanDetailProps {
  testId?: string;
  span?: TSpan;
  resultId?: string;
}

export interface TSpanDetailsComponentProps {
  span?: TSpan;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
}

const ComponentMap: Record<string, typeof Generic> = {
  [SemanticGroupNames.Http]: Http,
};

const getSpanTitle = (span: TSpan) => {
  const {primary, heading} = SpanService.getSpanNodeInfo(span);
  const spanTypeText = SemanticGroupNamesToText[span.type];

  return `${capitalize(heading) || spanTypeText} • ${primary} • ${span.name}`;
};

const SpanDetail: React.FC<TSpanDetailProps> = ({span}) => {
  const {open} = useAssertionForm();
  const Component = ComponentMap[span?.type || ''] || Generic;

  const title = (span && getSpanTitle(span)) || '';

  const onCreateAssertion = useCallback(
    ({value, key}: TSpanFlatAttribute) => {
      open({
        isEditing: false,
        defaultValues: {
          assertionList: [
            {
              comparator: CompareOperator.EQUALS,
              expected: value,
              attribute: key,
            },
          ],
          selectorList: [],
        },
      });
    },
    [open]
  );

  return (
    <S.SpanDetail>
      <SpanHeader title={title} />
      <Component span={span} onCreateAssertion={onCreateAssertion} />
    </S.SpanDetail>
  );
};

export default SpanDetail;
