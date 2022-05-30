import {Tooltip} from 'antd';
import {useCallback} from 'react';
import {useStore} from 'react-flow-renderer';

import AssertionCheckRow from 'components/AssertionCheckRow';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {clearAffectedSpans, setAffectedSpans} from 'redux/slices/TestDefinition.slice';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TAssertionResultEntry} from 'types/Assertion.types';
import * as S from './AssertionCard.styled';

interface TAssertionCardProps {
  assertionResult: TAssertionResultEntry;
  onSelectSpan(spanId: string): void;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
  selectedAssertion: string;
  setSelectedAssertion(assertion: string): void;
}

const AssertionCard: React.FC<TAssertionCardProps> = ({
  assertionResult: {selector, resultList, selectorList, pseudoSelector, spanIds},
  assertionResult,
  onSelectSpan,
  onDelete,
  onEdit,
  selectedAssertion,
  setSelectedAssertion,
}) => {
  const dispatch = useAppDispatch();
  const store = useStore();

  const spanCountText = `${spanIds.length} ${spanIds.length > 1 ? 'spans' : 'span'}`;
  const definition = useAppSelector(state => TestDefinitionSelectors.selectDefinitionBySelector(state, selector));
  const {isDraft = false, isDeleted = false} = definition || {};

  const getIsSelectedSpan = useCallback(
    (id: string): boolean => {
      const {selectedElements} = store.getState();
      const found = selectedElements ? selectedElements.find(element => element.id === id) : undefined;

      return Boolean(found);
    },
    [store]
  );

  const handleOnClick = () => {
    if (selectedAssertion === selector) {
      dispatch(clearAffectedSpans());
      setSelectedAssertion('');
      return;
    }
    setSelectedAssertion(selector);
    dispatch(setAffectedSpans(spanIds));
  };

  return (
    <S.AssertionCard data-cy="assertion-card" isSelected={selectedAssertion === selector}>
      <S.Header onClick={handleOnClick}>
        <div>
          <S.SelectorListText>
            {selectorList.map(({value}) => value).join(' ')} {pseudoSelector?.selector}
            {pseudoSelector?.number && `(${pseudoSelector?.number})`}
          </S.SelectorListText>
          <S.SpanCountText>{spanCountText}</S.SpanCountText>
        </div>
        <S.ActionsContainer>
          {isDraft && <S.StatusTag>draft</S.StatusTag>}
          {isDeleted && <S.StatusTag color="#61175E">deleted</S.StatusTag>}
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
              assertionSelectorList={selectorList.map(({value}) => value)}
            />
          ))
        )}
      </S.Body>
    </S.AssertionCard>
  );
};

export default AssertionCard;
