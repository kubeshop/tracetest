import {FC} from 'react';
import * as S from './Header.styled';

const Header: FC = () => {
  return (
    <S.Header>
      <S.TitleText>Tracetest</S.TitleText>
      <S.NavMenu>
        <S.NavMenuItem>GitHub</S.NavMenuItem>
        <S.NavMenuItem>Documentation</S.NavMenuItem>
      </S.NavMenu>
    </S.Header>
  );
};

export default Header;
