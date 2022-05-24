import {useCallback} from 'react';
import {useStore} from 'react-flow-renderer';
import {TAssertionResultEntry} from '../../types/Assertion.types';
import AssertionCheckRow from '../AssertionCheckRow';
import * as S from './AssertionCard.styled';

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

  const spanCountText = `${spanCount} ${spanCount > 1 ? 'spans' : 'span'}`;

  const getIsSelectedSpan = useCallback(
    (id: string): boolean => {
      const {selectedElements} = store.getState();
      const found = selectedElements ? selectedElements.find(element => element.id === id) : undefined;

      return Boolean(found);
    },
    [store]
  );

  return (
    <S.AssertionCard data-cy="assertion-card">
      <S.Header>
        <div>
          <S.SelectorListText>
            {selectorList.map(({value}) => value).join(' ')} {pseudoSelector?.selector}
            {pseudoSelector?.number && `(${pseudoSelector?.number})`}
          </S.SelectorListText>
          <S.SpanCountText>{spanCountText}</S.SpanCountText>
        </div>
        <div>
          <S.EditIcon data-cy="edit-assertion-button" onClick={() => onEdit(assertionResult)} />
          <S.DeleteIcon onClick={() => onDelete(selector)} />
        </div>
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
