import {Typography} from 'antd';
import {useCallback, useState} from 'react';
import {useTestDefinition} from '../../providers/TestDefinition/TestDefinition.provider';
import {TAssertionResultEntry, TAssertionResults} from '../../types/Assertion.types';
import AssertionCard from '../AssertionCard/AssertionCard';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import * as S from './AssertionCardList.styled';

interface TAssertionCardListProps {
  assertionResults: TAssertionResults;
  onSelectSpan(spanId: string): void;
  testId: string;
}

const AssertionCardList: React.FC<TAssertionCardListProps> = ({assertionResults: {resultList}, onSelectSpan}) => {
  const {open} = useAssertionForm();
  const {remove} = useTestDefinition();
  const [selectedAssertion, setSelectedAssertion] = useState('');

  const handleEdit = useCallback(
    ({selector, resultList: list, selectorList, pseudoSelector}: TAssertionResultEntry) => {
      open({
        isEditing: true,
        selector,
        defaultValues: {
          assertionList: list.map(({assertion}) => assertion),
          selectorList,
          pseudoSelector,
        },
      });
    },
    [open]
  );

  const handleDelete = useCallback(
    (selector: string) => {
      remove(selector);
    },
    [remove]
  );

  return (
    <S.AssertionCardList data-cy="assertion-card-list">
      {resultList.length ? (
        resultList.map(assertionResult =>
          assertionResult.resultList.length ? (
            <AssertionCard
              key={assertionResult.id}
              assertionResult={assertionResult}
              onSelectSpan={onSelectSpan}
              onEdit={handleEdit}
              onDelete={handleDelete}
              selectedAssertion={selectedAssertion}
              setSelectedAssertion={setSelectedAssertion}
            />
          ) : null
        )
      ) : (
        <S.EmptyStateContainer data-cy="empty-assertion-card-list">
          <S.EmptyStateIcon />
          <Typography.Text disabled>No Data</Typography.Text>
        </S.EmptyStateContainer>
      )}
    </S.AssertionCardList>
  );
};

export default AssertionCardList;
