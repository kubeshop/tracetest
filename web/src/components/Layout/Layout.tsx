import {FC} from 'react';
import useRouterSync from 'hooks/useRouterSync';
import ConfirmationModalProvider from 'providers/ConfirmationModal';
import * as S from './Layout.styled';
import Header from '../Header';
import FileViewerModalProvider from '../FileViewerModal/FileViewerModal.provider';

const Layout: FC = ({children}) => {
  useRouterSync();

  return (
    <FileViewerModalProvider>
      <ConfirmationModalProvider>
        <Header hasLogo />
        <S.Content>{children}</S.Content>
      </ConfirmationModalProvider>
    </FileViewerModalProvider>
  );
};

export default Layout;
