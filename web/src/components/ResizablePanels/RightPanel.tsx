import * as Spaces from 'react-spaces';
import useResizablePanel, {TPanel} from '../ResizablePanels/hooks/useResizablePanel';
import PanelContainer from './PanelContainer';
import * as S from './ResizablePanels.styled';

interface IProps {
  panel: TPanel;
  order?: number;
  children: React.ReactNode;
}

const RightPanel = ({panel, order = 1, children}: IProps) => {
  const {size, onChange, onStopResize} = useResizablePanel({panel});

  return (
    <Spaces.RightResizable
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
    </Spaces.RightResizable>
  );
};

export default RightPanel;
