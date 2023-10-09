import * as Spaces from 'react-spaces';
import {useLayoutEffect} from 'react';
import {noop} from 'lodash';
import useResizablePanel, {TPanel, TSize} from '../ResizablePanels/hooks/useResizablePanel';
import Splitter from '../ResizablePanels/Splitter';

interface IProps {
  panel: TPanel;
  order?: number;
  children(size: TSize): React.ReactNode;
  tooltip?: string;
  isToolTipVisible?: boolean;
  onOpen?(): void;
  dataTour?: string;
}

const LeftPanel = ({
  panel,
  order = 1,
  tooltip,
  isToolTipVisible = false,
  children,
  onOpen = noop,
  dataTour,
}: IProps) => {
  const {size, toggle, onStopResize} = useResizablePanel({panel});

  useLayoutEffect(() => {
    if (size.isOpen) onOpen();
  }, [onOpen, size.isOpen]);

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
          onClick={toggle}
          tooltip={tooltip}
          isToolTipVisible={isToolTipVisible}
          tooltipPlacement="right"
          dataTour={dataTour}
        />
      )}
    >
      {children(size)}
    </Spaces.LeftResizable>
  );
};

export default LeftPanel;
