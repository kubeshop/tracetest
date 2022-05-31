import {Tooltip} from 'antd';
import {useCallback} from 'react';
import {useStore} from 'react-flow-renderer';

import AssertionCheckRow from 'components/AssertionCheckRow';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TAssertionResultEntry} from 'types/Assertion.types';
import * as S from './AssertionCard.styled';

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
  const {setSelectedAssertion} = useTestDefinition();
  const store = useStore();

  const selectedAssertion = useAppSelector(TestDefinitionSelectors.selectSelectedAssertion);
  const definition = useAppSelector(state => TestDefinitionSelectors.selectDefinitionBySelector(state, selector));
  const {isDraft = false, isDeleted = false} = definition || {};
  const spanCountText = `${spanIds.length} ${spanIds.length > 1 ? 'spans' : 'span'}`;

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
      return setSelectedAssertion('');
    }
    setSelectedAssertion(selector);
  };

  return (
    <S.AssertionCard
      $isSelected={selectedAssertion === selector}
      data-cy="assertion-card"
      id={`assertion-${assertionResult.selector}`}
    >
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
