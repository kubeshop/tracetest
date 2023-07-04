import * as Spaces from 'react-spaces';
import useResizablePanel, {TPanel, TSize} from '../ResizablePanels/hooks/useResizablePanel';
import Splitter from '../ResizablePanels/Splitter';

interface IProps {
  panel: TPanel;
  order?: number;
  children(size: TSize): React.ReactNode;
  tooltip?: string;
}

const LeftPanel = ({panel, order = 1, tooltip, children}: IProps) => {
  const {size, toggle, onStopResize} = useResizablePanel({panel});

  return (
    <Spaces.LeftResizable
      onResizeEnd={newSize => onStopResize(newSize)}
      minimumSize={size.minSize}
      maximumSize={size.maxSize}
      size={size.size}
      key={size.name}
      order={order}
      handleRender={props => (
        <Splitter
          {...props}
          name={size.name}
          isOpen={size.isOpen}
          onClick={() => toggle()}
          tooltip={!size.isOpen ? tooltip : ''}
          tooltipPlacement="right"
        />
      )}
    >
      {children(size)}
    </Spaces.LeftResizable>
  );
};

export default LeftPanel;
