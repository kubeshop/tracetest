import {PanelResizeHandle} from 'react-resizable-panels';
import {TPanel} from '../ResizablePanels/hooks/useResizablePanel';
import Panel from './Panel';

interface IProps {
  panel: TPanel;
  children: React.ReactNode;
  onOpen?(): void;
}

const LeftPanel = (props: IProps) => {
  return (
    <>
      <PanelResizeHandle />
      <Panel {...props} />
    </>
  );
};

export default LeftPanel;
