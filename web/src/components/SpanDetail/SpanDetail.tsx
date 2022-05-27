import {capitalize} from 'lodash';
import {useCallback} from 'react';

import {useAssertionForm} from 'components/AssertionForm/AssertionFormProvider';
import {CompareOperator} from 'constants/Operator.constants';
import {SemanticGroupNames, SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import SpanService from 'services/Span.service';
import OperatorService from 'services/Operator.service';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpan, TSpanFlatAttribute} from 'types/Span.types';
import Generic from './components/Generic';
import Http from './components/Http';
import SpanHeader from './SpanHeader';
import * as S from './SpanDetail.styled';

export interface ISpanDetailsComponentProps {
  assertions?: TResultAssertions;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
  span?: TSpan;
}

interface IProps {
  span?: TSpan;
}

const ComponentMap: Record<string, typeof Generic> = {
  [SemanticGroupNames.Http]: Http,
};

const getSpanTitle = (span: TSpan) => {
  const {primary, heading} = SpanService.getSpanNodeInfo(span);
  const spanTypeText = SemanticGroupNamesToText[span.type];

  return `${capitalize(heading) || spanTypeText} • ${primary} • ${span.name}`;
};

const SpanDetail: React.FC<IProps> = ({span}) => {
  const {open} = useAssertionForm();
  const assertions = useAppSelector(state =>
    TestDefinitionSelectors.selectAssertionResultsBySpan(state, span?.id || '')
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
