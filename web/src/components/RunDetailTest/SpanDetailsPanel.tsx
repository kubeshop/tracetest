import {useSpan} from 'providers/Span/Span.provider';
import SpanDetail from '../SpanDetail/SpanDetail';
import {LeftPanel, PanelContainer} from '../ResizablePanels';

const panel = {
  name: 'SPAN_DETAILS',
  minSize: 25,
  maxSize: 320,
};

const SpanDetailsPanel = () => {
  const {selectedSpan} = useSpan();

  return (
    <LeftPanel
      panel={panel}
      tooltip="A certain span contains an attribute and this attribute has a specific value. You can check it here."
    >
      {size => (
        <PanelContainer $isOpen={size.isOpen}>
          <SpanDetail span={selectedSpan} />
        </PanelContainer>
      )}
    </LeftPanel>
  );
};

export default SpanDetailsPanel;
