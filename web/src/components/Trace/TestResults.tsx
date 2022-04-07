import {Typography} from 'antd';
import {FC, useMemo} from 'react';
import {runAssertionByTrace} from '../../services/AssertionService';
import {useGetTestByIdQuery} from '../../services/TestService';
import {ITrace} from '../../types';
import TraceAssertionsResultTable from '../TraceAssertionsTable/TraceAssertionsTable';
import * as S from './TestResults.styled';

type TTestResultsProps = {
  testId: string;
  trace: ITrace;
  onSpanSelected(spanId: string): void;
};

const TestResults: FC<TTestResultsProps> = ({testId, trace, onSpanSelected}) => {
  const {data: test} = useGetTestByIdQuery(testId);

  const traceResultList = useMemo(
    () => test?.assertions.map(assertion => runAssertionByTrace(trace, assertion)) || [],
    [test?.assertions, trace]
  );
  const totalSpanCount = trace.resourceSpans.length;
  const totalAssertionCount = test?.assertions.length;

  const [totalPassedCount, totalFailedCount] = useMemo(
    () =>
      traceResultList.reduce<[number, number]>(
        ([innerTotalPassedCount, innerTotalFailedCount], {spanListAssertionResult}) => {
          const [passed, failed] = spanListAssertionResult.reduce<[number, number]>(
            ([passedResultCount, failedResultCount], {resultList}) => {
              const passedCount = resultList.filter(({hasPassed}) => hasPassed).length;
              const failedCount = resultList.filter(({hasPassed}) => !hasPassed).length;

              return [passedResultCount + passedCount, failedResultCount + failedCount];
            },
            [0, 0]
          );

          return [innerTotalPassedCount + passed, innerTotalFailedCount + failed];
        },
        [0, 0]
      ),
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
      {traceResultList.length ? (
        traceResultList.map(assertionResult =>
          assertionResult.spanListAssertionResult.length ? (
            <TraceAssertionsResultTable
              key={assertionResult.assertion.assertionId}
              assertionResult={assertionResult}
              trace={trace}
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
    </S.Container>
  );
};

export default TestResults;
