import {capitalize} from 'lodash';
import {useCallback} from 'react';
import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import SpanService from 'services/Span.service';
import {TSpan, TSpanFlatAttribute} from 'types/Span.types';
import {CompareOperator} from 'constants/Operator.constants';
import OperatorService from 'services/Operator.service';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TResultAssertions} from 'types/Assertion.types';
import {useAssertionForm} from 'components/AssertionForm/AssertionForm.provider';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import SelectorService from 'services/Selector.service';
import { ResultViewModes } from 'constants/Test.constants';

import SpanDetailTabs from './SpanDetailTabs';
import SpanHeader from './SpanHeader';
import * as S from './SpanDetail.styled';
import {OPEN_BOTTOM_PANEL_STATE, useRunLayout} from '../RunLayout';

export interface ISpanDetailsComponentProps {
  assertions?: TResultAssertions;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
  span?: TSpan;
}

interface IProps {
  span?: TSpan;
}

const getSpanTitle = (span: TSpan) => {
  const {primary, heading} = SpanService.getSpanNodeInfo(span);
  const spanTypeText = SemanticGroupNamesToText[span.type];

  return `${capitalize(heading) || spanTypeText} • ${primary} • ${span.name}`;
};

const SpanDetail: React.FC<IProps> = ({span}) => {
  const {openBottomPanel} = useRunLayout();
  const {open} = useAssertionForm();
  const {viewResultsMode} = useTestDefinition();
  const assertions = useAppSelector(state =>
    TestDefinitionSelectors.selectAssertionResultsBySpan(state, span?.id || '')
  );
  const title = (span && getSpanTitle(span)) || '';

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
      <SpanHeader title={title} />
      <SpanDetailTabs onCreateAssertion={onCreateAssertion} span={span} assertions={assertions} />
    </S.SpanDetail>
  );
};

export default SpanDetail;
