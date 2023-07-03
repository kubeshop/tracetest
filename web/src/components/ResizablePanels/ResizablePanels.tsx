import * as Spaces from 'react-spaces';
import * as S from './ResizablePanels.styled';

const ResizablePanels: React.FC = ({children}) => {
  return (
    <>
      <S.GlobalStyle />
      <Spaces.Fixed height="100%" width="100vw">
        {children}
      </Spaces.Fixed>
    </>
  );
};

export default ResizablePanels;
