import {Panel as ResizablePanel} from 'react-resizable-panels';
import {useLayoutEffect} from 'react';
import {noop} from 'lodash';
import useResizablePanel, {TPanel} from './hooks/useResizablePanel';
import * as S from './ResizablePanels.styled';

interface IProps {
  panel: TPanel;
  children: React.ReactNode;
  onOpen?(): void;
}

const Panel = ({onOpen = noop, panel, children}: IProps) => {
  const {size, onStopResize, onChange, ref} = useResizablePanel({panel});

  useLayoutEffect(() => {
    if (size.isOpen) onOpen();
  }, [onOpen, size.isOpen]);

  return (
    <ResizablePanel ref={ref} maxSize={size.maxSize()} defaultSize={size.openSize()} onResize={onStopResize}>
      <S.PanelContainer $isOpen={size.isOpen} onClick={() => !size.isOpen && onChange(size.openSize())}>
        {children}
      </S.PanelContainer>
    </ResizablePanel>
  );
};

export default Panel;
