import {noop} from 'lodash';
import {useCallback} from 'react';

import {CompareOperator} from 'constants/Operator.constants';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {useAppSelector} from 'redux/hooks';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import OperatorService from 'services/Operator.service';
import SpanService from 'services/Span.service';
import {TSpan, TSpanFlatAttribute} from 'types/Span.types';
import Attributes from './Attributes';
import Header from './Header';
import * as S from './SpanDetail.styled';

interface IProps {
  onCreateTestSpec?(): void;
  searchText?: string;
  span?: TSpan;
}

const SpanDetail = ({onCreateTestSpec = noop, searchText, span}: IProps) => {
  const {open} = useTestSpecForm();
  const spansResult = useAppSelector(TestSpecsSelectors.selectSpansResult);
  const assertions = useAppSelector(state => TestSpecsSelectors.selectAssertionResultsBySpan(state, span?.id || ''));

  const handleCreateTestSpec = useCallback(
    ({value, key}: TSpanFlatAttribute) => {
      TraceAnalyticsService.onAddAssertionButtonClick();
      const selector = SpanService.getSelectorInformation(span!);

      open({
        isEditing: false,
        selector,
        defaultValues: {
          assertions: [
            {
              comparator: OperatorService.getOperatorSymbol(CompareOperator.EQUALS),
              expected: value,
              attribute: key,
            },
          ],
          selector,
        },
      });

      onCreateTestSpec();
    },
    [onCreateTestSpec, open, span]
  );

  return (
    <>
      <Header
        span={span}
        totalFailedChecks={span?.id ? spansResult[span.id]?.failed : 0}
        totalPassedChecks={span?.id ? spansResult[span?.id]?.passed : 0}
      />
      <S.HeaderDivider />
      <Attributes
        assertions={assertions}
        attributeList={span?.attributeList ?? []}
        onCreateTestSpec={handleCreateTestSpec}
        searchText={searchText}
        type={span?.type ?? SemanticGroupNames.General}
      />
    </>
  );
};

export default SpanDetail;
