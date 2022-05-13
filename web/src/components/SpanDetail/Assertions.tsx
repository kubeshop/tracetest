import {Typography} from 'antd';
import {PlusOutlined} from '@ant-design/icons';
import {ISpan} from 'types/Span.types';
import SkeletonTable from 'components/SkeletonTable';
import AssertionsResultTable from '../AssertionsTable';
import * as S from './SpanDetail.styled';
import {IAssertionResultList} from '../../types/Assertion.types';
import {useCreateAssertionModal} from '../CreateAssertionModal/CreateAssertionModalProvider';

interface IProps {
  span?: ISpan;
  testId?: string;
  resultId?: string;
  assertionsResultList: IAssertionResultList[];
}

const Assertion: React.FC<IProps> = ({assertionsResultList, testId, span, resultId}) => {
  const {open} = useCreateAssertionModal();

  return (
    <SkeletonTable loading={!span}>
      <S.AssertionActionsContainer>
        <S.AddAssertionButton
          data-cy="add-assertion-button"
          icon={<PlusOutlined />}
          onClick={() =>
            open({
              span,
              testId: testId!,
              resultId: resultId!,
            })
          }
        >
          Add Assertion
        </S.AddAssertionButton>
      </S.AssertionActionsContainer>
      {assertionsResultList.length ? (
        assertionsResultList
          .sort((a, b) => (a.assertion.assertionId > b.assertion.assertionId ? -1 : 1))
          .map(({assertion, assertionResultList}, index) => (
            <AssertionsResultTable
              key={assertion.assertionId}
              assertionResults={assertionResultList}
              sort={index + 1}
              assertion={assertion}
              span={span!}
              resultId={resultId!}
              testId={testId!}
            />
          ))
      ) : (
        <S.DetailsEmptyStateContainer data-cy="empty-assertion-table">
          <S.DetailsTableEmptyStateIcon />
          <Typography.Text disabled>Add assertion to see result here.</Typography.Text>
        </S.DetailsEmptyStateContainer>
      )}
    </SkeletonTable>
  );
};

export default Assertion;
