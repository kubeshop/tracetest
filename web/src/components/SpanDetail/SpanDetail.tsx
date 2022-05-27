import {capitalize} from 'lodash';
import {useCallback} from 'react';

import {useAssertionForm} from 'components/AssertionForm/AssertionFormProvider';
import {CompareOperator} from 'constants/Operator.constants';
import {SemanticGroupNames, SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {useAppSelector} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import SpanService from 'services/Span.service';
import OperatorService from 'services/Operator.service';
import {TSpan, TSpanFlatAttribute} from 'types/Span.types';
import Generic from './components/Generic';
import Http from './components/Http';
import SpanHeader from './SpanHeader';
import * as S from './SpanDetail.styled';

export interface ISpanDetailsComponentProps {
  assertions?: IResultAssertions;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
  span?: TSpan;
}

export interface IResultAssertions {
  [key: string]: {
    failed: IResult[];
    passed: IResult[];
  };
}

export interface IResult {
  id: string;
  label: string;
}

interface IProps {
  resultId?: string;
  span?: TSpan;
  testId?: string;
}

const ComponentMap: Record<string, typeof Generic> = {
  [SemanticGroupNames.Http]: Http,
};

const getSpanTitle = (span: TSpan) => {
  const {primary, heading} = SpanService.getSpanNodeInfo(span);
  const spanTypeText = SemanticGroupNamesToText[span.type];

  return `${capitalize(heading) || spanTypeText} • ${primary} • ${span.name}`;
};

const SpanDetail: React.FC<IProps> = ({span, testId, resultId}) => {
  const {open} = useAssertionForm();
  const assertions = useAppSelector(state =>
    AssertionSelectors.selectResultAssertionsBySpan(state, testId || '', resultId, span?.id || '')
  );

  const Component = ComponentMap[span?.type || ''] || Generic;
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
      <Component span={span} onCreateAssertion={onCreateAssertion} assertions={assertions} />
    </S.SpanDetail>
  );
};

export default SpanDetail;
