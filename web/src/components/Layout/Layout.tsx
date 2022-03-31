import {FC} from 'react';
import * as S from './Layout.styled';
import Header from '../Header';

const Layout: FC = ({children}) => {
  return (
    <>
      <Header />
      <S.Content>{children}</S.Content>
    </>
  );
};

export default Layout;
