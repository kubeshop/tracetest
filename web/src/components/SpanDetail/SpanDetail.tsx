import {useCallback} from 'react';

import {CompareOperator} from 'constants/Operator.constants';
import {useAssertionForm} from 'components/AssertionForm/AssertionForm.provider';
import {OPEN_BOTTOM_PANEL_STATE, useRunLayout} from 'components/RunLayout';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import OperatorService from 'services/Operator.service';
import SpanService from 'services/Span.service';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpan, TSpanFlatAttribute} from 'types/Span.types';
import * as S from './SpanDetail.styled';
import SpanDetailTabs from './SpanDetailTabs';
import SpanHeader from './SpanHeader';

export interface ISpanDetailComponentProps {
  assertions?: TResultAssertions;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
  span?: TSpan;
}

interface IProps {
  span?: TSpan;
}

const SpanDetail = ({span}: IProps) => {
  const {openBottomPanel} = useRunLayout();
  const {open} = useAssertionForm();
  const assertions = useAppSelector(state =>
    TestDefinitionSelectors.selectAssertionResultsBySpan(state, span?.id || '')
  );

  const onCreateAssertion = useCallback(
    ({value, key}: TSpanFlatAttribute) => {
      openBottomPanel(OPEN_BOTTOM_PANEL_STATE.FORM);
      TraceAnalyticsService.onAddAssertionButtonClick();
      const selector = SpanService.getSelectorInformation(span!);

      open({
        isEditing: false,
        selector,
        defaultValues: {
          assertionList: [
            {
              comparator: OperatorService.getOperatorSymbol(CompareOperator.EQUALS),
              expected: value,
              attribute: key,
            },
          ],
          selector,
        },
      });
    },
    [openBottomPanel, span, open]
  );

  return (
    <S.SpanDetail>
      <SpanHeader span={span} />
      <SpanDetailTabs assertions={assertions} onCreateAssertion={onCreateAssertion} span={span} />
    </S.SpanDetail>
  );
};

export default SpanDetail;
