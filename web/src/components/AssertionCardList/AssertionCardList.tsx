import {Typography} from 'antd';
import {useCallback} from 'react';
import {useDeleteAssertionMutation} from '../../redux/apis/Test.api';
import {IAssertionResult} from '../../types/Assertion.types';
import AssertionCard from '../AssertionCard/AssertionCard';
import {useCreateAssertionModal} from '../CreateAssertionModal/CreateAssertionModalProvider';
import * as S from './AssertionCardList.styled';

interface IAssertionCardListProps {
  assertionResultList: IAssertionResult[];
  onSelectSpan(spanId: string): void;
  resultId: string;
  testId: string;
}

const AssertionCardList: React.FC<IAssertionCardListProps> = ({
  assertionResultList,
  onSelectSpan,
  testId,
  resultId,
}) => {
  const {open} = useCreateAssertionModal();
  const [deleteAssertion] = useDeleteAssertionMutation();

  const handleEdit = useCallback(
    ({assertion, spanListAssertionResult}: IAssertionResult) => {
      const [{span}] = spanListAssertionResult;

      open({
        span,
        testId,
        assertion,
        resultId,
      });
    },
    [open, resultId, testId]
  );

  return (
    <S.AssertionCardList>
      {assertionResultList.length ? (
        assertionResultList.map(assertionResult =>
          assertionResult.spanListAssertionResult.length ? (
            <AssertionCard
              key={assertionResult.assertion?.assertionId}
              assertionResult={assertionResult}
              onSelectSpan={onSelectSpan}
              onEdit={handleEdit}
              onDelete={assertionId => deleteAssertion({testId, assertionId})}
            />
          ) : null
        )
      ) : (
        <S.EmptyStateContainer>
          <S.EmptyStateIcon />
          <Typography.Text disabled>No Data</Typography.Text>
        </S.EmptyStateContainer>
      )}
    </S.AssertionCardList>
  );
};

export default AssertionCardList;
