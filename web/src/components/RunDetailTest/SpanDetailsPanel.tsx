import {useSpan} from 'providers/Span/Span.provider';
import SpanDetail from '../SpanDetail/SpanDetail';
import {LeftPanel, PanelContainer} from '../ResizablePanels';

const panel = {
  name: 'SPAN_DETAILS',
  minSize: 15,
  maxSize: 320,
};

const SpanDetailsPanel = () => {
  const {selectedSpan} = useSpan();

  return (
    <LeftPanel panel={panel}>
      {size => (
        <PanelContainer $isOpen={size.isOpen}>
          <SpanDetail span={selectedSpan} />
        </PanelContainer>
      )}
    </LeftPanel>
  );
};

export default SpanDetailsPanel;
