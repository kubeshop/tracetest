import {PlusOutlined} from '@ant-design/icons';
import {Badge} from 'antd';
import * as React from 'react';
import {MouseEventHandler, useCallback, useMemo} from 'react';
import {TTestRun} from 'types/TestRun.types';
import Date from 'utils/Date';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import SpanService from '../../services/Span.service';
import TraceService from '../../services/Trace.service';
import {TAssertionResults} from '../../types/Assertion.types';
import {TSpan} from '../../types/Span.types';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {Steps} from '../GuidedTour/traceStepList';
import {useSetIsCollapsedCallbackDirection} from '../ResizableDrawer/useSetIsCollapsedCallback';
import * as S from './TraceDrawer.styled';
import {useChevronDirectionMemo} from './useChevronDirectionMemo';

interface IProps {
  visiblePortion: number;
  run: TTestRun;
  assertionResults?: TAssertionResults;
  isDisabled: boolean;
  selectedSpan: TSpan;
  height?: number;
  min?: number;
  max?: number;
}

const TraceDrawerHeader: React.FC<IProps> = ({
  run: {createdAt},
  visiblePortion,
  assertionResults,
  isDisabled,
  selectedSpan,
  height,
  max,
  min,
}) => {
  const $isCollapsed = useChevronDirectionMemo(height, max, min);
  const onClick = useSetIsCollapsedCallbackDirection($isCollapsed);
  const {open} = useAssertionForm();
  const totalAssertionCount = assertionResults?.resultList.length || 0;

  const {totalPassedCount, totalFailedCount} = useMemo(
    () => TraceService.getTestResultCount(assertionResults!),
    [assertionResults]
  );

  const handleAssertionClick: MouseEventHandler<HTMLElement> = useCallback(
    event => {
      event.stopPropagation();
      const {selectorList, pseudoSelector} = SpanService.getSelectorInformation(selectedSpan!);

      open({
        isEditing: false,
        defaultValues: {
          pseudoSelector,
          selectorList,
        },
      });
    },
    [open, selectedSpan]
  );
  return (
    <S.Header
      visiblePortion={visiblePortion}
      data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Timeline)}
      style={{height: visiblePortion, minHeight: visiblePortion}}
      onClick={onClick}
    >
      <div>
        <S.HeaderText strong>Test Result</S.HeaderText>
        <S.StartDateText>{Date.format(createdAt)}</S.StartDateText>
        <S.HeaderText strong>
          {totalAssertionCount} assertion(s) • {totalPassedCount + totalFailedCount} check(s) •{' '}
          <Badge count="P" style={{backgroundColor: '#49AA19'}} /> <S.CountNumber>{totalPassedCount}</S.CountNumber>
          <Badge count="F" /> <S.CountNumber>{totalFailedCount}</S.CountNumber>
        </S.HeaderText>
      </div>
      <div>
        <S.AddAssertionButton
          data-cy="add-assertion-button"
          icon={<PlusOutlined />}
          disabled={isDisabled}
          onClick={handleAssertionClick}
        >
          Add Assertion
        </S.AddAssertionButton>
        <span style={{marginLeft: 16}}>
          <S.Chevron $isCollapsed={$isCollapsed} />
        </span>
      </div>
    </S.Header>
  );
};

export default TraceDrawerHeader;
