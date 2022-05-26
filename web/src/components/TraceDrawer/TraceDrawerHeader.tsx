import {PlusOutlined} from '@ant-design/icons';
import {Badge} from 'antd';
import {format, parseISO} from 'date-fns';
import {useMemo} from 'react';
import {TTestRun} from 'types/TestRun.types';
import {CompareOperator} from '../../constants/Operator.constants';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import OperatorService from '../../services/Operator.service';
import TraceService from '../../services/Trace.service';
import {TAssertionResults} from '../../types/Assertion.types';
import {TSpan} from '../../types/Span.types';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {Steps} from '../GuidedTour/traceStepList';
import * as S from './TraceDrawer.styled';

interface IProps {
  visiblePortion: number;
  run: TTestRun;
  assertionResults?: TAssertionResults;
  onClick(): void;
  isDisabled: boolean;
  selectedSpan: TSpan;
}

const TraceDrawerHeader: React.FC<IProps> = ({
  run: {trace, createdAt},
  visiblePortion,
  assertionResults,
  onClick,
  isDisabled,
  selectedSpan,
}) => {
  const {open} = useAssertionForm();
  const totalSpanCount = trace?.spans.length;
  const totalAssertionCount = assertionResults?.resultList.length || 0;

  const {totalPassedCount, totalFailedCount} = useMemo(
    () => TraceService.getTestResultCount(assertionResults!),
    [assertionResults]
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
        <S.HeaderText strong>Test Results</S.HeaderText>
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
            open({
              isEditing: false,
              defaultValues: {
                selectorList:
                  selectedSpan?.signature.map(attribute => ({
                    value: attribute.value,
                    key: attribute.key,
                    operator: OperatorService.getOperatorSymbol(CompareOperator.EQUALS),
                  })) || [],
              },
            });
          }}
        >
          Add Assertion
        </S.AddAssertionButton>
      </div>
    </S.Header>
  );
};

export default TraceDrawerHeader;
