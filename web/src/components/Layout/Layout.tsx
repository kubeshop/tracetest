import {FC} from 'react';
import * as S from './Layout.styled';
import Header from '../Header';
import FileViewerModalProvider from '../FileViewerModal/FileViewerModal.provider';

const Layout: FC = ({children}) => {
  return (
    <FileViewerModalProvider>
      <Header />
      <S.Content>{children}</S.Content>
    </FileViewerModalProvider>
  );
};

export default Layout;
