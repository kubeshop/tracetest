import {Typography} from 'antd';
import {FC, useMemo} from 'react';
import {getTestResultCount} from '../../services/TraceService';
import {AssertionResult, ITrace} from '../../types';
import TraceAssertionsResultTable from '../TraceAssertionsTable/TraceAssertionsTable';
import * as S from './TestResults.styled';

type TTestResultsProps = {
  trace: ITrace;
  traceResultList: AssertionResult[];
  onSpanSelected(spanId: string): void;
};

const TestResults: FC<TTestResultsProps> = ({trace, traceResultList, onSpanSelected}) => {
  const totalSpanCount = trace.resourceSpans.length;
  const totalAssertionCount = traceResultList.length || 0;

  const {totalPassedCount, totalFailedCount} = useMemo(() => getTestResultCount(traceResultList), [traceResultList]);

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
