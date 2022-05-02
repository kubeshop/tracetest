import {Typography} from 'antd';
import {useSelector} from 'react-redux';
import SkeletonTable from 'components/SkeletonTable';
import {FC, useMemo} from 'react';
import {ITrace} from '../../types/Trace.types';
import TraceAssertionsResultTable from '../TraceAssertionsTable/TraceAssertionsTable';
import TraceService from '../../services/Trace.service';
import * as S from './TestResults.styled';
import TestResultSelectors from '../../selectors/TestResult.selectors';

type TTestResultsProps = {
  trace?: ITrace;
  resultId: string;
  onSpanSelected(spanId: string): void;
};

const TestResults: FC<TTestResultsProps> = ({trace, resultId, onSpanSelected}) => {
  const traceResultList = useSelector(TestResultSelectors.selectTestResultList(resultId));
  const totalSpanCount = trace?.spans.length;
  const totalAssertionCount = traceResultList.length || 0;

  const {totalPassedCount, totalFailedCount} = useMemo(
    () => TraceService.getTestResultCount(traceResultList),
    [traceResultList]
  );

  return (
    <S.Container>
      <S.Header>
        <Typography.Text strong>
          {totalSpanCount} total spans • {totalAssertionCount} selectors • {totalPassedCount + totalFailedCount} checks
          • {totalPassedCount} passed/{totalFailedCount} failed
        </Typography.Text>
      </S.Header>
      <SkeletonTable loading={!trace}>
        {traceResultList.length ? (
          traceResultList.map(assertionResult =>
            assertionResult.spanListAssertionResult.length ? (
              <TraceAssertionsResultTable
                key={assertionResult.assertion.assertionId}
                assertionResult={assertionResult}
                onSpanSelected={onSpanSelected}
              />
            ) : null
          )
        ) : (
          <S.EmptyStateContainer>
            <S.EmptyStateIcon />
            <Typography.Text disabled>No Data</Typography.Text>
          </S.EmptyStateContainer>
        )}
      </SkeletonTable>
    </S.Container>
  );
};

export default TestResults;
