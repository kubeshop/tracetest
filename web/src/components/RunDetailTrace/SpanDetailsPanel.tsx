import {useNavigate} from 'react-router-dom';
import {useCallback} from 'react';
import {useAppSelector} from 'redux/hooks';
import TraceSelectors from 'selectors/Trace.selectors';
import TestRun from 'models/TestRun.model';
import SpanSelectors from 'selectors/Span.selectors';
import SpanDetail from '../SpanDetail/SpanDetail';
import {LeftPanel, PanelContainer} from '../ResizablePanels';

interface IProps {
  run: TestRun;
  testId: string;
}

const panel = {
  name: 'SPAN_DETAILS',
  minSize: 15,
  maxSize: 320,
};

const SpanDetailsPanel = ({run, testId}: IProps) => {
  const searchText = useAppSelector(TraceSelectors.selectSearchText);
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);
  const navigate = useNavigate();
  const span = useAppSelector(state => SpanSelectors.selectSpanById(state, selectedSpan, testId, run.id));

  const handleOnCreateSpec = useCallback(() => {
    navigate(`/test/${testId}/run/${run.id}/test`);
  }, [navigate, run.id, testId]);

  return (
    <LeftPanel panel={panel}>
      {size => (
        <PanelContainer $isOpen={size.isOpen}>
          <SpanDetail onCreateTestSpec={handleOnCreateSpec} searchText={searchText} span={span} />
        </PanelContainer>
      )}
    </LeftPanel>
  );
};

export default SpanDetailsPanel;
