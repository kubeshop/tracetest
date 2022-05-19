import {PlusOutlined} from '@ant-design/icons';
import {Badge} from 'antd';
import {format, parseISO} from 'date-fns';
import {useMemo} from 'react';
import {useSelector} from 'react-redux';
import {ITestRunResult} from 'types/TestRunResult.types';
import TestResultSelectors from '../../selectors/TestResult.selectors';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import TraceService from '../../services/Trace.service';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {Steps} from '../GuidedTour/traceStepList';
import * as S from './TraceDrawer.styled';

interface IProps {
  visiblePortion: number;
  result: ITestRunResult;
  onClick(): void;
  isDisabled: boolean;
}

const TraceDrawerHeader: React.FC<IProps> = ({
  result: {resultId, trace, createdAt},
  visiblePortion,
  onClick,
  isDisabled,
}) => {
  const {open} = useAssertionForm();
  const traceResultList = useSelector(TestResultSelectors.selectTestResultList(resultId));
  const totalSpanCount = trace?.spans.length;
  const totalAssertionCount = traceResultList.length || 0;

  const {totalPassedCount, totalFailedCount} = useMemo(
    () => TraceService.getTestResultCount(traceResultList),
    [traceResultList]
  );

  const startDate = useMemo(() => format(parseISO(createdAt), "EEEE, do MMMM yyyy 'at' HH:mm:ss"), [createdAt]);

  return (
    <S.Header
      visiblePortion={visiblePortion}
      data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Timeline)}
      style={{height: visiblePortion}}
      onClick={onClick}
    >
      <div>
        <S.HeaderText strong>Trace Information</S.HeaderText>
        <S.StartDateText>Trace Start {startDate}</S.StartDateText>
        <S.HeaderText strong>
          {totalSpanCount} total span(s) • {totalAssertionCount} assertion(s) • {totalPassedCount + totalFailedCount}{' '}
          check(s) • <Badge count="P" style={{backgroundColor: '#49AA19'}} />{' '}
          <S.CountNumber>{totalPassedCount}</S.CountNumber>
          <Badge count="F" /> <S.CountNumber>{totalFailedCount}</S.CountNumber>
        </S.HeaderText>
      </div>
      <div>
        <S.AddAssertionButton
          data-cy="add-assertion-button"
          icon={<PlusOutlined />}
          disabled={isDisabled}
          onClick={event => {
            event.stopPropagation();
            open();
          }}
        >
          Add Assertion
        </S.AddAssertionButton>
      </div>
    </S.Header>
  );
};

export default TraceDrawerHeader;
