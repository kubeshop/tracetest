import {useCallback} from 'react';
import {useAppSelector} from 'redux/hooks';
import SpanDetail, {TraceAttributeRow, TraceSubHeader} from 'components/SpanDetail';
import TestRun from 'models/TestRun.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import SpanSelectors from 'selectors/Span.selectors';
import TraceSelectors from 'selectors/Trace.selectors';
import useAttributePanelTooltip from 'hooks/useAttributePanelTooltip';
import {LeftPanel, PanelContainer} from '../ResizablePanels';

interface IProps {
  run: TestRun;
  testId: string;
}

const panel = {
  name: 'SPAN_DETAILS',
  minSize: 25,
  maxSize: 320,
  isDefaultOpen: true,
};

const SpanDetailsPanel = ({run, testId}: IProps) => {
  const searchText = useAppSelector(TraceSelectors.selectSearchText);
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);
  const {navigate} = useDashboard();
  const span = useAppSelector(state => SpanSelectors.selectSpanById(state, selectedSpan, testId, run.id));
  const {onClose, tooltip, isVisible} = useAttributePanelTooltip();

  const handleOnCreateSpec = useCallback(() => {
    navigate(`/test/${testId}/run/${run.id}/test`);
  }, [navigate, run.id, testId]);

  return (
    <LeftPanel panel={panel} tooltip={tooltip} onOpen={onClose} isToolTipVisible={isVisible}>
      {size => (
        <PanelContainer $isOpen={size.isOpen}>
          <SpanDetail
            onCreateTestSpec={handleOnCreateSpec}
            searchText={searchText}
            span={span}
            AttributeRowComponent={TraceAttributeRow}
            SubHeaderComponent={TraceSubHeader}
          />
        </PanelContainer>
      )}
    </LeftPanel>
  );
};

export default SpanDetailsPanel;
