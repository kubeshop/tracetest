import {Panel as ResizablePanel} from 'react-resizable-panels';
import {useCallback, useLayoutEffect} from 'react';
import {noop} from 'lodash';
import useResizablePanel, {TPanel} from './hooks/useResizablePanel';
import * as S from './ResizablePanels.styled';
import Handle from './Handle';

interface IProps {
  panel: TPanel;
  children: React.ReactNode;
  onOpen?(): void;
  handlePlacement?: 'left' | 'right';
}

const Panel = ({onOpen = noop, panel, children, handlePlacement = 'right'}: IProps) => {
  const {size, onStopResize, onChange, ref} = useResizablePanel({panel});

  useLayoutEffect(() => {
    if (size.isOpen) onOpen();
  }, [onOpen, size.isOpen]);

  const toggle = useCallback(() => {
    onChange(size.isOpen ? size.closeSize() : size.openSize());
  }, [onChange, size]);

  return (
    <>
      {handlePlacement === 'left' && <Handle placement={handlePlacement} isOpen={!size.isOpen} onToggle={toggle} />}
      <ResizablePanel ref={ref} maxSize={size.maxSize()} defaultSize={size.openSize()} onResize={onStopResize}>
        <S.PanelContainer $isOpen={size.isOpen} onClick={() => !size.isOpen && onChange(size.openSize())}>
          {children}
        </S.PanelContainer>
      </ResizablePanel>
      {handlePlacement === 'right' && <Handle placement={handlePlacement} isOpen={size.isOpen} onToggle={toggle} />}
    </>
  );
};

export default Panel;
