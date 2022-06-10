import {PlusOutlined} from '@ant-design/icons';
import {Badge} from 'antd';
import {MouseEventHandler, useCallback, useMemo} from 'react';

import {useAssertionForm} from 'components/AssertionForm/AssertionFormProvider';
import {Steps} from 'components/GuidedTour/traceStepList';
import {OPEN_BOTTOM_PANEL_STATE, useRunLayout} from 'components/RunLayout';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import SpanService from 'services/Span.service';
import TraceService from 'services/Trace.service';
import {TAssertionResults} from 'types/Assertion.types';
import {TSpan} from 'types/Span.types';
import {TTestRun} from 'types/TestRun.types';
import Date from 'utils/Date';
import * as S from './RunBottomPanel.styled';

interface IProps {
  run: TTestRun;
  assertionResults?: TAssertionResults;
  isDisabled: boolean;
  selectedSpan: TSpan;
}

const Header: React.FC<IProps> = ({run: {createdAt}, assertionResults, isDisabled, selectedSpan}) => {
  const {isBottomPanelOpen, openBottomPanel, toggleBottomPanel} = useRunLayout();
  const {open} = useAssertionForm();
  const totalAssertionCount = assertionResults?.resultList.length || 0;

  const {totalPassedCount, totalFailedCount} = useMemo(
    () => TraceService.getTestResultCount(assertionResults!),
    [assertionResults]
  );

  const handleAssertionClick: MouseEventHandler<HTMLElement> = useCallback(
    event => {
      event.stopPropagation();
      openBottomPanel(OPEN_BOTTOM_PANEL_STATE.FORM);
      const {selectorList, pseudoSelector} = SpanService.getSelectorInformation(selectedSpan!);

      open({
        isEditing: false,
        defaultValues: {
          pseudoSelector,
          selectorList,
        },
      });
    },
    [open, selectedSpan, openBottomPanel]
  );
  return (
    <S.Header
      data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Timeline)}
      onClick={() => toggleBottomPanel()}
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
        <S.ChevronContainer>
          <S.Chevron $isCollapsed={isBottomPanelOpen} />
        </S.ChevronContainer>
      </div>
    </S.Header>
  );
};

export default Header;
