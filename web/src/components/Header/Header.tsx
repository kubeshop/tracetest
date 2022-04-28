import {FC} from 'react';
import {Link, useLocation} from 'react-router-dom';
import {DOCUMENTATION_URL, GITHUB_URL} from '../../constants/Common.contants';
import * as S from './Header.styled';
import Logo from '../../assets/Logo.svg';

const Header: FC = () => {
  const {pathname} = useLocation();

  return (
    <S.Header>
      <Link to="/">
        <S.TitleText>
          <img src={Logo} />
        </S.TitleText>
      </Link>
      <S.NavMenu selectedKeys={[pathname]}>
        <S.NavMenuItem key={GITHUB_URL}>
          <a href={GITHUB_URL} target="_blank">
            GitHub
          </a>
        </S.NavMenuItem>
        <S.NavMenuItem key={DOCUMENTATION_URL}>
          <a href={DOCUMENTATION_URL} target="_blank">
            Documentation
          </a>
        </S.NavMenuItem>
      </S.NavMenu>
    </S.Header>
  );
};

export default Header;
