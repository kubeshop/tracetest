import {FC} from 'react';
import {Link, useLocation} from 'react-router-dom';
import * as S from './Header.styled';

const Header: FC = () => {
  const {pathname} = useLocation();

  return (
    <S.Header>
      <Link to="/">
        <S.TitleText>Tracetest</S.TitleText>
      </Link>
      <S.NavMenu selectedKeys={[pathname]}>
        <S.NavMenuItem>
          <a href="https://github.com/kubeshop/tracetest" target="_blank">
            GitHub
          </a>
        </S.NavMenuItem>
        <S.NavMenuItem>
          <a href="https://kubeshop.github.io/tracetest/" target="_blank">
            Documentation
          </a>
        </S.NavMenuItem>
      </S.NavMenu>
    </S.Header>
  );
};

export default Header;
