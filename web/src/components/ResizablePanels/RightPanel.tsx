import {TPanel} from '../ResizablePanels/hooks/useResizablePanel';
import Panel from './Panel';

interface IProps {
  panel: TPanel;
  children: React.ReactNode;
  onOpen?(): void;
}

const RightPanel = (props: IProps) => {
  return <Panel handlePlacement="left" {...props} />;
};

export default RightPanel;
