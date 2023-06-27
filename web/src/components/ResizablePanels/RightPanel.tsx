import * as Spaces from 'react-spaces';
import Splitter from '../ResizablePanels/Splitter';
import useResizablePanel, {TPanel, TSize} from '../ResizablePanels/hooks/useResizablePanel';

interface IProps {
  panel: TPanel;
  order?: number;
  children(size: TSize): React.ReactNode;
}

const RightPanel = ({panel, order = 1, children}: IProps) => {
  const {size, toggle, onStopResize} = useResizablePanel({panel});

  return (
    <Spaces.RightResizable
      onResizeEnd={newSize => onStopResize(newSize)}
      minimumSize={size.minSize}
      maximumSize={size.maxSize}
      size={size.size}
      key={size.name}
      order={order}
      handleRender={props => <Splitter {...props} name={size.name} isOpen={!size.isOpen} onClick={() => toggle()} />}
    >
      {children(size)}
    </Spaces.RightResizable>
  );
};

export default RightPanel;
