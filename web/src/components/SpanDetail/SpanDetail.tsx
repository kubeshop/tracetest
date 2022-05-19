import {capitalize} from 'lodash';
import {useCallback} from 'react';
import {SemanticGroupNames, SemanticGroupNamesToText} from '../../constants/SemanticGroupNames.constants';
import SpanService from '../../services/Span.service';
import {ISpan, ISpanFlatAttribute} from '../../types/Span.types';
import Generic from './components/Generic';
import Http from './components/Http';
import SpanHeader from './SpanHeader';
import * as S from './SpanDetail.styled';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {CompareOperator} from '../../constants/Operator.constants';

export interface ISpanDetailProps {
  testId?: string;
  span?: ISpan;
  resultId?: string;
}

export interface ISpanDetailsComponentProps {
  span?: ISpan;
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

const SpanDetail: React.FC<ISpanDetailProps> = ({span}) => {
  const {open} = useAssertionForm();
  const Component = ComponentMap[span?.type || ''] || Generic;

  const title = (span && getSpanTitle(span)) || '';

  const onCreateAssertion = useCallback(
    ({type, value, key}: ISpanFlatAttribute) => {
      open({
        isEditing: false,
        defaultValues: {
          assertionList: [
            {
              compareOp: CompareOperator.EQUALS,
              value,
              key,
              type,
            },
          ],
          selectorList: span?.signature || [],
        },
      });
    },
    [open, span?.signature]
  );

  return (
    <S.SpanDetail>
      <SpanHeader title={title} />
      <Component span={span} onCreateAssertion={onCreateAssertion} />
    </S.SpanDetail>
  );
};

export default SpanDetail;
