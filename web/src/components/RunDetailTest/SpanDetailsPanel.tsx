import {StepsID} from 'components/GuidedTour/testRunSteps';
import SpanDetail, {TestAttributeRow, TestSubHeader} from 'components/SpanDetail';
import {useSpan} from 'providers/Span/Span.provider';
import useAttributePanelTooltip from 'hooks/useAttributePanelTooltip';
import {LeftPanel, PanelContainer} from '../ResizablePanels';

const panel = {
  name: 'SPAN_DETAILS',
  minSize: 25,
  maxSize: 320,
  isDefaultOpen: true,
};

const SpanDetailsPanel = () => {
  const {selectedSpan} = useSpan();
  const {tooltip, isVisible, onClose} = useAttributePanelTooltip();

  return (
    <LeftPanel
      panel={panel}
      tooltip={tooltip}
      isToolTipVisible={isVisible}
      onOpen={onClose}
      dataTour={StepsID.SpanDetails}
    >
      {size => (
        <PanelContainer $isOpen={size.isOpen}>
          <SpanDetail span={selectedSpan} AttributeRowComponent={TestAttributeRow} SubHeaderComponent={TestSubHeader} />
        </PanelContainer>
      )}
    </LeftPanel>
  );
};

export default SpanDetailsPanel;
