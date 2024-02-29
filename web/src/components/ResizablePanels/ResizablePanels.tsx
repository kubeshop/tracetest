import {PanelGroup} from 'react-resizable-panels';
import * as S from './ResizablePanels.styled';

interface IProps {
  saveId?: string;
}

const ResizablePanels: React.FC<IProps> = ({children, saveId}) => {
  return (
    <>
      <S.GlobalStyle />
      <PanelGroup autoSaveId={saveId} direction="horizontal">
        {children}
      </PanelGroup>
    </>
  );
};

export default ResizablePanels;
