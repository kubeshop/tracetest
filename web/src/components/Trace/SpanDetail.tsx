import {Button, Typography} from 'antd';
import {PlusOutlined} from '@ant-design/icons';
import {FC, useState} from 'react';
import AssertionsResultTable from 'components/AssertionsTable/AssertionsTable';
import CreateAssertionModal from 'components/CreateAssertionModal';
import SkeletonTable from 'components/SkeletonTable';
import * as S from './SpanDetail.styled';
import Attributes from './Attributes';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';
import {useAppSelector} from '../../redux/hooks';
import AssertionSelectors from '../../selectors/Assertion.selectors';
import {ISpan} from '../../types/Span.types';

const {onAddAssertionButtonClick} = TraceAnalyticsService;

type TSpanDetailProps = {
  testId?: string;
  targetSpan?: ISpan;
  resultId?: string;
};

const SpanDetail: FC<TSpanDetailProps> = ({testId, targetSpan, resultId}) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);

  const assertionsResultList = useAppSelector(
    AssertionSelectors.selectAssertionResultListBySpan(testId, resultId, targetSpan?.spanId)
  );

  return (
    <>
      <S.DetailsContainer>
        <S.DetailsHeader>
          <Typography.Text strong> Details for selected span</Typography.Text>
          <Button
            type="link"
            icon={<PlusOutlined />}
            onClick={() => {
              onAddAssertionButtonClick();
              setOpenCreateAssertion(true);
            }}
          >
            Add Assertion
          </Button>
        </S.DetailsHeader>
        <SkeletonTable loading={!targetSpan}>
          {assertionsResultList.length ? (
            assertionsResultList
              .sort((a, b) => (a.assertion.assertionId > b.assertion.assertionId ? -1 : 1))
              .map(({assertion, assertionResultList}, index) => (
                <AssertionsResultTable
                  key={assertion.assertionId}
                  assertionResults={assertionResultList}
                  sort={index + 1}
                  assertion={assertion}
                  span={targetSpan!}
                  resultId={resultId!}
                  testId={testId!}
                />
              ))
          ) : (
            <S.DetailsEmptyStateContainer>
              <S.DetailsTableEmptyStateIcon />
              <Typography.Text disabled>No Data</Typography.Text>
            </S.DetailsEmptyStateContainer>
          )}
        </SkeletonTable>
      </S.DetailsContainer>
      <Attributes spanAttributeList={targetSpan?.attributeList} />

      {targetSpan?.spanId && testId && (
        <CreateAssertionModal
          key={`KEY_${targetSpan?.spanId}`}
          testId={testId}
          span={targetSpan!}
          resultId={resultId!}
          open={openCreateAssertion}
          onClose={() => setOpenCreateAssertion(false)}
        />
      )}
    </>
  );
};

export default SpanDetail;
