import {Tooltip} from 'antd';
import {useCallback} from 'react';
import {useStore} from 'react-flow-renderer';
import {useAppSelector} from '../../redux/hooks';
import TestDefinitionSelectors from '../../selectors/TestDefinition.selectors';
import {TAssertionResultEntry} from '../../types/Assertion.types';
import AssertionCheckRow from '../AssertionCheckRow';
import * as S from './AssertionCard.styled';
import AssertionCardSelectorList from './AssertionCardSelectorList';

interface TAssertionCardProps {
  assertionResult: TAssertionResultEntry;
  onSelectSpan(spanId: string): void;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
}

const AssertionCard: React.FC<TAssertionCardProps> = ({
  assertionResult: {selector, resultList, selectorList, pseudoSelector, spanCount},
  assertionResult,
  onSelectSpan,
  onDelete,
  onEdit,
}) => {
  const store = useStore();
  const {selectedElements} = store.getState();

  const spanCountText = `${spanCount} ${spanCount > 1 ? 'spans' : 'span'}`;
  const {isDraft = false, isDeleted = false} =
    useAppSelector(state => TestDefinitionSelectors.selectDefinitionBySelector(state, selector)) || {};

  const getIsSelectedSpan = useCallback(
    (id: string): boolean => {
      const found = selectedElements ? selectedElements.find(element => element.id === id) : undefined;

      return Boolean(found);
    },
    [selectedElements]
  );

  return (
    <S.AssertionCard data-cy="assertion-card" id={`assertion-${assertionResult.selector}`}>
      <S.Header>
        <AssertionCardSelectorList selectorList={selectorList} pseudoSelector={pseudoSelector} />
        <S.ActionsContainer>
          {isDraft && <S.StatusTag>draft</S.StatusTag>}
          {isDeleted && <S.StatusTag color="#61175E">deleted</S.StatusTag>}
          <S.SpanCountText>{spanCountText}</S.SpanCountText>
          <Tooltip color="white" title="Edit Assertion">
            <S.EditIcon data-cy="edit-assertion-button" onClick={() => onEdit(assertionResult)} />
          </Tooltip>
          <Tooltip color="white" title="Delete Assertion">
            <S.DeleteIcon onClick={() => onDelete(selector)} />
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
