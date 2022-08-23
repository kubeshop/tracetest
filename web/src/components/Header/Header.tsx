import {Link} from 'react-router-dom';

import Logo from 'assets/Logo.svg';
import * as S from './Header.styled';
import {HeaderMenu} from './HeaderMenu';

const Header = () => (
  <S.Header>
    <Link to="/">
      <S.Logo alt="tracetest logo" data-cy="logo" src={Logo} />
    </Link>
    <HeaderMenu />
  </S.Header>
);

export default Header;
