import {Tooltip} from 'antd';
import AssertionCheckRow from 'components/AssertionCheckRow';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useCallback} from 'react';
import {useStore} from 'react-flow-renderer';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {useAppSelector} from '../../redux/hooks';
import AssertionAnalyticsService from '../../services/Analytics/AssertionAnalytics.service';
import {TAssertionResultEntry} from '../../types/Assertion.types';
import * as S from './AssertionCard.styled';
import AssertionCardSelectorList from './AssertionCardSelectorList';

interface TAssertionCardProps {
  assertionResult: TAssertionResultEntry;
  onSelectSpan(spanId: string): void;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
}

const AssertionCard: React.FC<TAssertionCardProps> = ({
  assertionResult: {selector, resultList, selectorList, pseudoSelector, spanIds},
  assertionResult,
  onSelectSpan,
  onDelete,
  onEdit,
}) => {
  const {setSelectedAssertion, revert} = useTestDefinition();
  const store = useStore();
  const {selectedElements} = store.getState();

  const selectedAssertion = useAppSelector(TestDefinitionSelectors.selectSelectedAssertion);
  const {
    isDraft = false,
    isDeleted = false,
    originalSelector = '',
  } = useAppSelector(state => TestDefinitionSelectors.selectDefinitionBySelector(state, selector)) || {};
  const spanCountText = `${spanIds.length} ${spanIds.length > 1 ? 'spans' : 'span'}`;

  const getIsSelectedSpan = useCallback(
    (id: string): boolean => {
      const found = selectedElements ? selectedElements.find(element => element.id === id) : undefined;

      return Boolean(found);
    },
    [selectedElements]
  );

  const handleOnClick = () => {
    if (selectedAssertion === selector) {
      return setSelectedAssertion('');
    }

    AssertionAnalyticsService.onAssertionClick();
    setSelectedAssertion(selector);
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
          <AssertionCardSelectorList selectorList={selectorList} pseudoSelector={pseudoSelector} />
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
