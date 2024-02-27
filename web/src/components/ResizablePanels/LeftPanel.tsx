import * as Spaces from 'react-spaces';
import {useLayoutEffect} from 'react';
import {noop} from 'lodash';
import useResizablePanel, {TPanel} from '../ResizablePanels/hooks/useResizablePanel';
import PanelContainer from './PanelContainer';
import * as S from './ResizablePanels.styled';

interface IProps {
  panel: TPanel;
  order?: number;
  children: React.ReactNode;
  onOpen?(): void;
}

const LeftPanel = ({panel, order = 1, children, onOpen = noop}: IProps) => {
  const {size, onStopResize, onChange} = useResizablePanel({panel});

  useLayoutEffect(() => {
    if (size.isOpen) onOpen();
  }, [onOpen, size.isOpen]);

  return (
    <Spaces.LeftResizable
      allowOverflow
      onResizeEnd={newSize => onStopResize(newSize)}
      minimumSize={size.minSize()}
      maximumSize={size.maxSize()}
      size={size.size}
      key={size.name}
      order={order}
      handleRender={props => <S.SplitterContainer {...props} />}
    >
      <PanelContainer size={size} onChange={onChange}>
        {children}
      </PanelContainer>
    </Spaces.LeftResizable>
  );
};

export default LeftPanel;
