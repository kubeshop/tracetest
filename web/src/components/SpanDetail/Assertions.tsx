import {useState} from 'react';
import {Typography} from 'antd';
import {PlusOutlined} from '@ant-design/icons';
import {ISpan} from 'types/Span.types';
import SkeletonTable from 'components/SkeletonTable';
import AssertionsResultTable from '../AssertionsTable';
import * as S from './SpanDetail.styled';
import CreateAssertionModal from '../CreateAssertionModal';
import {IAssertionResultList} from '../../types/Assertion.types';

interface IProps {
  span?: ISpan;
  testId?: string;
  resultId?: string;
  assertionsResultList: IAssertionResultList[];
}

const Assertion: React.FC<IProps> = ({assertionsResultList, testId, span, resultId}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <SkeletonTable loading={!span}>
      <S.AssertionActionsContainer>
        <S.AddAssertionButton
          data-cy="add-assertion-button"
          icon={<PlusOutlined />}
          onClick={() => setIsModalOpen(true)}
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
          <Typography.Text disabled>No Data</Typography.Text>
        </S.DetailsEmptyStateContainer>
      )}

      <CreateAssertionModal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        span={span!}
        testId={testId!}
        resultId={resultId!}
      />
    </SkeletonTable>
  );
};

export default Assertion;
