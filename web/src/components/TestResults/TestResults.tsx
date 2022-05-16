import {Badge} from 'antd';
import {useSelector} from 'react-redux';
import {format, parseISO} from 'date-fns';
import SkeletonTable from 'components/SkeletonTable';
import {FC, PointerEventHandler, useMemo} from 'react';
import * as S from './TestResults.styled';
import TestResultSelectors from '../../selectors/TestResult.selectors';
import TraceService from '../../services/Trace.service';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/traceStepList';
import {ITestRunResult} from '../../types/TestRunResult.types';
import AssertionCardList from '../AssertionCardList';

type TTestResultsProps = {
  result: ITestRunResult;
  onPointerDown?: PointerEventHandler;
  visiblePortion: number;
  onSpanSelected(spanId: string): void;
  onHeaderClick(): void;
};

const TestResults: FC<TTestResultsProps> = ({
  result: {resultId, testId, trace, createdAt},
  onSpanSelected,
  visiblePortion,
  onPointerDown,
  onHeaderClick,
}) => {
  const traceResultList = useSelector(TestResultSelectors.selectTestResultList(resultId));
  const totalSpanCount = trace?.spans.length;
  const totalAssertionCount = traceResultList.length || 0;

  const {totalPassedCount, totalFailedCount} = useMemo(
    () => TraceService.getTestResultCount(traceResultList),
    [traceResultList]
  );

  const startDate = format(parseISO(createdAt), "EEEE, do MMMM yyyy 'at' HH:mm:ss");

  return (
    <S.Container data-cy="test-results">
      <S.Header
        onPointerDown={onPointerDown}
        visiblePortion={visiblePortion}
        data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Timeline)}
        style={{height: visiblePortion}}
        onClick={onHeaderClick}
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
          <AssertionCardList
            assertionResultList={traceResultList}
            onSelectSpan={onSpanSelected}
            resultId={resultId}
            testId={testId}
          />
        </SkeletonTable>
      </S.Content>
    </S.Container>
  );
};

export default TestResults;
