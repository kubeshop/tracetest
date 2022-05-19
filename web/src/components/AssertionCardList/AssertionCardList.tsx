import {Typography} from 'antd';
import {useCallback} from 'react';
import {useDeleteAssertionMutation} from '../../redux/apis/Test.api';
import AssertionService from '../../services/Assertion.service';
import {IAssertionResult} from '../../types/Assertion.types';
import AssertionCard from '../AssertionCard/AssertionCard';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import * as S from './AssertionCardList.styled';

interface IAssertionCardListProps {
  assertionResultList: IAssertionResult[];
  onSelectSpan(spanId: string): void;
  resultId: string;
  testId: string;
}

const AssertionCardList: React.FC<IAssertionCardListProps> = ({assertionResultList, onSelectSpan, testId}) => {
  const {open} = useAssertionForm();
  const [deleteAssertion] = useDeleteAssertionMutation();

  const handleEdit = useCallback(
    ({assertion: {assertionId, selectors: selectorList, spanAssertions}}: IAssertionResult) => {
      open({
        isEditing: true,
        assertionId,
        defaultValues: {
          selectorList,
          assertionList: AssertionService.parseSelectorSpanToAssertionSpan(spanAssertions),
        },
      });
    },
    [open]
  );

  return (
    <S.AssertionCardList data-cy="assertion-card-list">
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
        <S.EmptyStateContainer data-cy="empty-assertion-card-list">
          <S.EmptyStateIcon />
          <Typography.Text disabled>No Data</Typography.Text>
        </S.EmptyStateContainer>
      )}
    </S.AssertionCardList>
  );
};

export default AssertionCardList;
