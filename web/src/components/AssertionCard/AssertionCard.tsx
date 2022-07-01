import {Tooltip} from 'antd';
import {useCallback} from 'react';
import {useAppSelector} from 'redux/hooks';
import AssertionCheckRow from 'components/AssertionCheckRow';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import {TAssertionResultEntry} from 'types/Assertion.types';
import * as S from './AssertionCard.styled';
import AssertionCardSelectorList from './AssertionCardSelectorList';

interface IProps {
  assertionResult: TAssertionResultEntry;
  onSelectSpan(spanId: string): void;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
}

const AssertionCard = ({
  assertionResult: {selector, resultList, selectorList, pseudoSelector, spanIds, isAdvancedSelector},
  assertionResult,
  onSelectSpan,
  onDelete,
  onEdit,
}: IProps) => {
  const {setSelectedAssertion, revert, viewResultsMode} = useTestDefinition();
  const {onSetFocusedSpan, selectedSpan} = useSpan();
  const selectedAssertion = useAppSelector(TestDefinitionSelectors.selectSelectedAssertion);
  const {
    isDraft = false,
    isDeleted = false,
    originalSelector = '',
  } = useAppSelector(state => TestDefinitionSelectors.selectDefinitionBySelector(state, selector)) || {};
  const spanCountText = `${spanIds.length} ${spanIds.length > 1 ? 'spans' : 'span'}`;

  const getIsSelectedSpan = useCallback((id: string) => selectedSpan?.id === id, [selectedSpan]);

  const handleOnClick = () => {
    onSetFocusedSpan('');
    if (selectedAssertion === selector) {
      return setSelectedAssertion();
    }

    AssertionAnalyticsService.onAssertionClick();
    setSelectedAssertion(assertionResult);
  };

  const resetDefinition: React.MouseEventHandler = useCallback(
    e => {
      e.stopPropagation();
      AssertionAnalyticsService.onRevertAssertion();
      revert(originalSelector);
    },
    [revert, originalSelector]
  );

  return (
    <S.AssertionCard
      $isSelected={selectedAssertion === selector}
      data-cy="assertion-card"
      id={`assertion-${assertionResult.selector}`}
    >
      <S.Header onClick={handleOnClick}>
        <div>
          <AssertionCardSelectorList
            viewResultsMode={viewResultsMode}
            isAdvancedSelector={isAdvancedSelector}
            selector={selector}
            selectorList={selectorList}
            pseudoSelector={pseudoSelector}
          />
        </div>
        <S.ActionsContainer>
          {isDraft && <S.StatusTag>draft</S.StatusTag>}
          {isDeleted && <S.StatusTag color="#61175E">deleted</S.StatusTag>}
          <S.SpanCountText>{spanCountText}</S.SpanCountText>
          {isDraft && (
            <Tooltip color="white" title="Revert Assertion">
              <S.UndoIcon data-cy="assertion-action-revert" onClick={resetDefinition} />
            </Tooltip>
          )}
          <Tooltip color="white" title="Edit Assertion">
            <S.EditIcon
              data-cy="edit-assertion-button"
              onClick={e => {
                e.stopPropagation();
                onEdit(assertionResult);
              }}
            />
          </Tooltip>
          <Tooltip color="white" title="Delete Assertion">
            <S.DeleteIcon
              onClick={e => {
                e.stopPropagation();
                onDelete(selector);
              }}
            />
          </Tooltip>
        </S.ActionsContainer>
      </S.Header>
      <S.Body>
        {resultList.flatMap(({spanResults, assertion: {attribute}, assertion}) =>
          spanResults.map(result => (
            <AssertionCheckRow
              key={`${attribute}-${result.spanId}`}
              assertion={assertion}
              result={result}
              onSelectSpan={onSelectSpan}
              getIsSelectedSpan={getIsSelectedSpan}
            />
          ))
        )}
      </S.Body>
    </S.AssertionCard>
  );
};

export default AssertionCard;
