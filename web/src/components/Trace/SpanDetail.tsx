import {Button, Typography} from 'antd';
import {skipToken} from '@reduxjs/toolkit/dist/query';
import {PlusOutlined} from '@ant-design/icons';
import {FC, useMemo, useState} from 'react';
import {useGetTestByIdQuery} from 'services/TestService';
import {Assertion, ISpan, ITrace, SpanAssertionResult} from 'types';
import {runAssertionBySpanId} from 'services/AssertionService';
import AssertionsResultTable from 'components/AssertionsTable/AssertionsTable';
import SkeletonTable from 'components/SkeletonTable';
import * as S from './SpanDetail.styled';
import CreateAssertionModal from './CreateAssertionModal';
import Attributes from './Attributes';

type TSpanDetailProps = {
  testId: string;
  targetSpan?: ISpan;
  trace?: ITrace;
};

const SpanDetail: FC<TSpanDetailProps> = ({testId, targetSpan, trace}) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);
  const {data: test} = useGetTestByIdQuery(testId ?? skipToken);

  const assertionsResultList = useMemo(() => {
    if (!targetSpan?.spanId || !trace) {
      return [];
    }
    return (
      test?.assertions?.reduce<Array<{assertion: Assertion; assertionResultList: Array<SpanAssertionResult>}>>(
        (resultList, assertion) => {
          const assertionResultList = runAssertionBySpanId(targetSpan.spanId, trace, assertion);

          return assertionResultList ? [...resultList, {assertion, assertionResultList}] : resultList;
        },
        []
      ) || []
    );
  }, [targetSpan?.spanId, test?.assertions, trace]);

  return (
    <>
      <S.DetailsContainer>
        <S.DetailsHeader>
          <Typography.Text strong> Details for selected span</Typography.Text>
          <Button type="link" icon={<PlusOutlined />} onClick={() => setOpenCreateAssertion(true)}>
            Add Assertion
          </Button>
        </S.DetailsHeader>
        <SkeletonTable loading={!targetSpan || !trace}>
          {assertionsResultList.length ? (
            assertionsResultList.map(({assertion, assertionResultList}, index) => (
              <AssertionsResultTable
                key={assertion.assertionId}
                assertionResults={assertionResultList}
                sort={index + 1}
                assertion={assertion}
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
      <Attributes spanId={targetSpan?.spanId} trace={trace} />

      {targetSpan?.spanId && (
        <CreateAssertionModal
          key={`KEY_${targetSpan?.spanId}`}
          testId={testId}
          trace={trace}
          span={targetSpan!}
          open={openCreateAssertion}
          onClose={() => setOpenCreateAssertion(false)}
        />
      )}
    </>
  );
};

export default SpanDetail;
