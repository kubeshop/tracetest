import {Typography} from 'antd';
import {useCallback} from 'react';
import {TAssertionResults} from '../../types/Assertion.types';
import AssertionCard from '../AssertionCard/AssertionCard';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import * as S from './AssertionCardList.styled';

interface TAssertionCardListProps {
  assertionResults: TAssertionResults;
  onSelectSpan(spanId: string): void;
  testId: string;
}

const AssertionCardList: React.FC<TAssertionCardListProps> = ({
  assertionResults: {resultList},
  onSelectSpan,
  testId,
}) => {
  const {open} = useAssertionForm();

  const handleEdit = useCallback(data => {
    console.log('@@onEditAssertion', data);
  }, []);

  return (
    <S.AssertionCardList data-cy="assertion-card-list">
      {resultList.length ? (
        resultList.map(assertionResult =>
          assertionResult.resultList.length ? (
            <AssertionCard
              key={assertionResult.id}
              assertionResult={assertionResult}
              onSelectSpan={onSelectSpan}
              selectorList={[]}
              onEdit={handleEdit}
              onDelete={assertionId => console.log('@@onDeleteAssertion', assertionId)}
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
