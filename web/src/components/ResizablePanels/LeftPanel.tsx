import {TPanel} from '../ResizablePanels/hooks/useResizablePanel';
import Panel from './Panel';

interface IProps {
  panel: TPanel;
  children: React.ReactNode;
  onOpen?(): void;
}

const LeftPanel = (props: IProps) => {
  return <Panel {...props} />;
};

export default LeftPanel;
