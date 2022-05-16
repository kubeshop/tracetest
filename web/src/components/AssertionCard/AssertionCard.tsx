import {useCallback, useMemo} from 'react';
import {useStore} from 'react-flow-renderer';
import {IAssertionResult} from '../../types/Assertion.types';
import AssertionCheckRow from '../AssertionCheckRow';
import * as S from './AssertionCard.styled';

interface IAssertionCardProps {
  assertionResult: IAssertionResult;
  onSelectSpan(spanId: string): void;
  onDelete(assertionId: string): void;
  onEdit(assertionResult: IAssertionResult): void;
}

const AssertionCard: React.FC<IAssertionCardProps> = ({
  assertionResult: {assertion: {selectors = [], assertionId = ''} = {}, spanListAssertionResult},
  assertionResult,
  onSelectSpan,
  onDelete,
  onEdit,
}) => {
  const store = useStore();

  const spanCountText = useMemo(() => {
    const spanCount = spanListAssertionResult.length;

    return `${spanCount} ${spanCount > 1 ? 'spans' : 'span'}`;
  }, [spanListAssertionResult.length]);

  const getIsSelectedSpan = useCallback(
    (id: string): boolean => {
      const {selectedElements} = store.getState();
      const found = selectedElements ? selectedElements.find(element => element.id === id) : undefined;

      return Boolean(found);
    },
    [store]
  );

  return (
    <S.AssertionCard>
      <S.Header>
        <div>
          <S.SelectorListText>{selectors.map(({value}) => value).join(' ')}</S.SelectorListText>
          <S.SpanCountText>{spanCountText}</S.SpanCountText>
        </div>
        <div>
          <S.EditIcon onClick={() => onEdit(assertionResult)} />
          <S.DeleteIcon onClick={() => onDelete(assertionId)} />
        </div>
      </S.Header>
      <S.Body>
        {spanListAssertionResult.flatMap(({span, resultList}) =>
          resultList.map(result => (
            <AssertionCheckRow
              key={`${result.propertyName}-${span.spanId}`}
              result={result}
              span={span}
              onSelectSpan={onSelectSpan}
              getIsSelectedSpan={getIsSelectedSpan}
              assertionSelectorList={selectors.map(({value}) => value)}
            />
          ))
        )}
      </S.Body>
    </S.AssertionCard>
  );
};

export default AssertionCard;
