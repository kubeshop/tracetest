import {Typography} from 'antd';
import SkeletonTable from 'components/SkeletonTable';
import {FC, useMemo} from 'react';
import { TAssertionResult } from '../../entities/Assertion/Assertion.types';
import {getTestResultCount} from '../../entities/Trace/Trace.service';
import { TTrace } from '../../entities/Trace/Trace.types';
import TraceAssertionsResultTable from '../TraceAssertionsTable/TraceAssertionsTable';
import * as S from './TestResults.styled';

type TTestResultsProps = {
  trace?: TTrace;
  traceResultList: TAssertionResult[];
  onSpanSelected(spanId: string): void;
};

const TestResults: FC<TTestResultsProps> = ({trace, traceResultList, onSpanSelected}) => {
  const totalSpanCount = trace?.resourceSpans?.length;
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
      <SkeletonTable loading={!trace}>
        {traceResultList.length ? (
          traceResultList.map(assertionResult =>
            assertionResult.spanListAssertionResult.length ? (
              <TraceAssertionsResultTable
                key={assertionResult.assertion.assertionId}
                assertionResult={assertionResult}
                trace={trace!}
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
