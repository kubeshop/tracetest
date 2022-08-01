import {useCallback} from 'react';

import {CompareOperator} from 'constants/Operator.constants';
import {ResultViewModes} from 'constants/Test.constants';
import {useAssertionForm} from 'components/AssertionForm/AssertionForm.provider';
import {OPEN_BOTTOM_PANEL_STATE, useRunLayout} from 'components/RunLayout';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import OperatorService from 'services/Operator.service';
import SelectorService from 'services/Selector.service';
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
  const {viewResultsMode} = useTestDefinition();
  const spansResult = useAppSelector(TestDefinitionSelectors.selectSpansResult);
  const assertions = useAppSelector(state =>
    TestDefinitionSelectors.selectAssertionResultsBySpan(state, span?.id || '')
  );

  const onCreateAssertion = useCallback(
    ({value, key}: TSpanFlatAttribute) => {
      openBottomPanel(OPEN_BOTTOM_PANEL_STATE.FORM);
      TraceAnalyticsService.onAddAssertionButtonClick();
      const {selectorList, pseudoSelector} = SpanService.getSelectorInformation(span!);
      const selector = SelectorService.getSelectorString(selectorList, pseudoSelector);

      open({
        isEditing: false,
        selector,
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
          selector,
          isAdvancedSelector: viewResultsMode === ResultViewModes.Advanced,
        },
      });
    },
    [openBottomPanel, span, open, viewResultsMode]
  );

  return (
    <S.SpanDetail>
      <SpanHeader
        span={span}
        totalFailedChecks={span?.id ? spansResult[span.id]?.failed : 0}
        totalPassedChecks={span?.id ? spansResult[span?.id]?.passed : 0}
      />
      <SpanDetailTabs assertions={assertions} onCreateAssertion={onCreateAssertion} span={span} />
    </S.SpanDetail>
  );
};

export default SpanDetail;
