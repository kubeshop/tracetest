import {Button, Typography} from 'antd';
import {PlusOutlined} from '@ant-design/icons';
import {FC, useState} from 'react';
import AssertionsResultTable from 'components/AssertionsTable/AssertionsTable';
import CreateAssertionModal from 'components/CreateAssertionModal';
import SkeletonTable from 'components/SkeletonTable';
import * as S from '../SpanDetail.styled';
import {Attributes} from '../../Trace/TraceComponent/Attributes';
import TraceAnalyticsService from '../../../services/Analytics/TraceAnalytics.service';
import {useAppSelector} from '../../../redux/hooks';
import AssertionSelectors from '../../../selectors/Assertion.selectors';
import {ISpanDetailProps} from '../SpanDetail';

const {onAddAssertionButtonClick} = TraceAnalyticsService;

const GenericSpanDetail: FC<ISpanDetailProps> = ({testId, span, resultId}) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);

  const assertionsResultList = useAppSelector(
    AssertionSelectors.selectAssertionResultListBySpan(testId, resultId, span?.spanId)
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
        <SkeletonTable loading={!span}>
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
            <S.DetailsEmptyStateContainer>
              <S.DetailsTableEmptyStateIcon />
              <Typography.Text disabled>No Data</Typography.Text>
            </S.DetailsEmptyStateContainer>
          )}
        </SkeletonTable>
      </S.DetailsContainer>
      <Attributes spanAttributeList={span?.attributeList} />

      {span && testId && (
        <CreateAssertionModal
          key={`KEY_${span?.spanId}`}
          testId={testId}
          span={span}
          resultId={resultId!}
          open={openCreateAssertion}
          onClose={() => setOpenCreateAssertion(false)}
        />
      )}
    </>
  );
};

export default GenericSpanDetail;
