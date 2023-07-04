import {useNavigate} from 'react-router-dom';
import {useCallback} from 'react';
import {useAppSelector} from 'redux/hooks';
import SpanDetail, {TraceAttributeRow, TraceSubHeader} from 'components/SpanDetail';
import TestRun from 'models/TestRun.model';
import SpanSelectors from 'selectors/Span.selectors';
import TraceSelectors from 'selectors/Trace.selectors';
import {LeftPanel, PanelContainer} from '../ResizablePanels';

interface IProps {
  run: TestRun;
  testId: string;
}

const panel = {
  name: 'SPAN_DETAILS',
  minSize: 25,
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
    <LeftPanel
      panel={panel}
      tooltip="A certain span contains an attribute and this attribute has a specific value. You can check it here."
    >
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
