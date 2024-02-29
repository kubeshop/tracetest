import SpanDetail, {TestAttributeRow, TestSubHeader} from 'components/SpanDetail';
import {useSpan} from 'providers/Span/Span.provider';
import useAttributePanelTooltip from 'hooks/useAttributePanelTooltip';
import {LeftPanel} from '../ResizablePanels';

const panel = {
  isDefaultOpen: true,
  openSize: () => (window.innerWidth / 4 / window.innerWidth) * 100,
};

const SpanDetailsPanel = () => {
  const {selectedSpan} = useSpan();
  const {onClose} = useAttributePanelTooltip();

  return (
    <LeftPanel panel={panel} onOpen={onClose}>
      <SpanDetail span={selectedSpan} AttributeRowComponent={TestAttributeRow} SubHeaderComponent={TestSubHeader} />
    </LeftPanel>
  );
};

export default SpanDetailsPanel;
