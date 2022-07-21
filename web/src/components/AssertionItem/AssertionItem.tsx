import {useMemo} from 'react';

import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {ResultViewModes} from 'constants/Test.constants';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import AssertionService from 'services/Assertion.service';
import {TAssertionResultEntry} from 'types/Assertion.types';
import {TTrace} from 'types/Trace.types';
import AssertionActions from './AssertionActions';
import AssertionHeader from './AssertionHeader';
import * as S from './AssertionItem.styled';
import CheckItem from './CheckItem';
import SpanHeader from './SpanHeader';
import SpanActions from './SpanActions';

interface IProps {
  assertionResult: TAssertionResultEntry;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
  onRevertAssertion(originalSelector: string): void;
  onSelectSpan(spanId: string): void;
  onSetFocusedSpan(spanId: string): void;
  onSetSelectedAssertion(assertionResult?: TAssertionResultEntry): void;
  selectedAssertion: string;
  selectedSpan?: string;
  trace?: TTrace;
  viewResultsMode: ResultViewModes;
}

const AssertionItem = ({
  assertionResult: {isAdvancedSelector, pseudoSelector, resultList, selector, selectorList, spanIds},
  assertionResult,
  onDelete,
  onEdit,
  onRevertAssertion,
  onSelectSpan,
  onSetFocusedSpan,
  onSetSelectedAssertion,
  selectedAssertion,
  selectedSpan,
  trace,
  viewResultsMode,
}: IProps) => {
  const {
    isDeleted = false,
    isDraft = false,
    originalSelector = '',
  } = useAppSelector(state => TestDefinitionSelectors.selectDefinitionBySelector(state, selector)) || {};
  const totalPassedChecks = useMemo(() => AssertionService.getTotalPassedChecks(resultList), [resultList]);
  const results = useMemo(() => AssertionService.getResultsHashedBySpanId(resultList), [resultList]);

  return (
    <div data-cy="assertion-card" id={`assertion-${selector}`}>
      <S.AssertionCollapse expandIconPosition="right" $isSelected={selector === selectedAssertion}>
        <S.AssertionCollapse.Panel
          extra={
            <AssertionActions
              isDeleted={isDeleted}
              isDraft={isDraft}
              isSelected={selectedAssertion === selector}
              onDelete={() => {
                onDelete(selector);
              }}
              onEdit={() => {
                onEdit(assertionResult);
              }}
              onRevert={() => {
                AssertionAnalyticsService.onRevertAssertion();
                onRevertAssertion(originalSelector);
              }}
              onSelect={() => {
                onSetFocusedSpan('');
                if (selectedAssertion === selector) {
                  return onSetSelectedAssertion();
                }

                AssertionAnalyticsService.onAssertionClick();
                onSetSelectedAssertion(assertionResult);
              }}
            />
          }
          header={
            <AssertionHeader
              affectedSpans={spanIds.length}
              failedChecks={totalPassedChecks?.false ?? 0}
              isAdvancedMode={viewResultsMode === ResultViewModes.Advanced}
              isAdvancedSelector={isAdvancedSelector}
              passedChecks={totalPassedChecks?.true ?? 0}
              pseudoSelector={pseudoSelector}
              selectorList={selectorList}
              title={selector}
            />
          }
          key={`assertion-${selector}`}
        >
          {Object.entries(results).map(([spanId, checkResults]) => {
            const span = trace?.spans.find(({id}) => id === spanId);

            return (
              <S.SpanCard
                extra={<SpanActions isValid={Boolean(span)} onSelect={() => onSelectSpan(spanId)} />}
                key={`${assertionResult.id}-${spanId}`}
                size="small"
                title={<SpanHeader span={span} />}
                type="inner"
                $isSelected={spanId === selectedSpan}
                $type={span?.type ?? SemanticGroupNames.General}
              >
                {checkResults.map(checkResult => (
                  <CheckItem
                    check={checkResult}
                    key={`${checkResult.result.spanId}-${checkResult.assertion.attribute}`}
                  />
                ))}
              </S.SpanCard>
            );
          })}
        </S.AssertionCollapse.Panel>
      </S.AssertionCollapse>
    </div>
  );
};

export default AssertionItem;
