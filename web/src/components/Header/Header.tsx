import Logo from 'assets/Logo.svg';
import {FC} from 'react';
import {Link, useLocation} from 'react-router-dom';
import HomeAnalyticsService from '../../services/Analytics/HomeAnalytics.service';
import * as S from './Header.styled';
import {HeaderMenu} from './HeaderMenu';
import './index.css';

const {onGuidedTourClick} = HomeAnalyticsService;

const Header: FC = () => {
  const {pathname} = useLocation();
  return (
    <S.Header>
      <Link to="/">
        <S.TitleText>
          <img alt="tracetest_log" data-cy="logo" src={Logo} />
        </S.TitleText>
      </Link>
      <HeaderMenu onGuidedTourClick={onGuidedTourClick} pathname={pathname} />
    </S.Header>
  );
};

export default Header;
