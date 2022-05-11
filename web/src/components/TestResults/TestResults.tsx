import {Badge, Typography} from 'antd';
import {useSelector} from 'react-redux';
import {format, parseISO} from 'date-fns';
import SkeletonTable from 'components/SkeletonTable';
import {FC, Dispatch, PointerEventHandler, useMemo, SetStateAction, useEffect} from 'react';
import * as S from './TestResults.styled';
import TestResultSelectors from '../../selectors/TestResult.selectors';
import TraceService from '../../services/Trace.service';
import TraceAssertionsResultTable from '../TraceAssertionsTable/TraceAssertionsTable';
import {useElementSize} from '../../hooks/useElementSize';
import {useDoubleClick} from '../../hooks/useDoubleClick';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/traceStepList';
import {ITestRunResult} from '../../types/TestRunResult.types';

type TTestResultsProps = {
  result: ITestRunResult;
  onPointerDown?: PointerEventHandler;
  visiblePortion: number;
  setMax: Dispatch<SetStateAction<number>>;
  max: number;
  height?: number;
  setHeight?: Dispatch<SetStateAction<number>>;
  onSpanSelected(spanId: string): void;
};

const TestResults: FC<TTestResultsProps> = ({
  result: {resultId, trace, createdAt},
  onSpanSelected,
  setMax,
  setHeight,
  visiblePortion,
  onPointerDown,
  ...props
}) => {
  const [squareRef, {height}] = useElementSize();
  useEffect(() => setMax(height), [height, setMax]);

  const traceResultList = useSelector(TestResultSelectors.selectTestResultList(resultId));
  const totalSpanCount = trace?.spans.length;
  const totalAssertionCount = traceResultList.length || 0;

  const {totalPassedCount, totalFailedCount} = useMemo(
    () => TraceService.getTestResultCount(traceResultList),
    [traceResultList]
  );

  const startDate = format(parseISO(createdAt), "EEEE, do MMMM yyyy 'at' HH:mm:ss");

  return (
    <S.Container
      data-cy="test-results"
      ref={squareRef}
      onClick={useDoubleClick(() => setHeight?.(height === props.height ? visiblePortion : height))}
    >
      <S.Header
        onPointerDown={onPointerDown}
        visiblePortion={visiblePortion}
        data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Timeline)}
        style={{height: visiblePortion}}
      >
        <S.HeaderText strong>Trace Information</S.HeaderText>
        <S.StartDateText>Trace Start {startDate}</S.StartDateText>
        <S.HeaderText strong>
          {totalSpanCount} total spans • {totalAssertionCount} selectors • {totalPassedCount + totalFailedCount} checks
          • <Badge count="P" style={{backgroundColor: '#49AA19'}} /> <S.CountNumber>{totalPassedCount}</S.CountNumber>
          <Badge count="F" /> <S.CountNumber>{totalFailedCount}</S.CountNumber>
        </S.HeaderText>
      </S.Header>
      <S.Content>
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
      </S.Content>
    </S.Container>
  );
};

export default TestResults;
