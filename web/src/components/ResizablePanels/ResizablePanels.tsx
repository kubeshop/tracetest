// import * as Spaces from 'react-spaces';
import {PanelGroup} from 'react-resizable-panels';
import * as S from './ResizablePanels.styled';

const ResizablePanels: React.FC = ({children}) => {
  return (
    <>
      <S.GlobalStyle />
      <PanelGroup direction="horizontal">
        {children}
      </PanelGroup>
    </>
  );
};

export default ResizablePanels;
