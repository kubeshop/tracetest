import {FC} from 'react';
import {Menu} from 'antd';
import {Link, useLocation} from 'react-router-dom';
import {QuestionCircleOutlined} from '@ant-design/icons';
import {DOCUMENTATION_URL, GITHUB_URL} from 'constants/Common.constants';
import {useTour} from '@reactour/tour';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import Logo from 'assets/Logo.svg';
import * as S from './Header.styled';

const {onGuidedTourClick} = HomeAnalyticsService;

const Header: FC = () => {
  const {pathname} = useLocation();
  const {setIsOpen} = useTour();

  const handleGuidedTourCLick = () => {
    setIsOpen(true);
    onGuidedTourClick();
  };

  return (
    <S.Header>
      <Link to="/">
        <S.TitleText>
          <img data-cy="logo" src={Logo} />
        </S.TitleText>
      </Link>
      <S.NavMenu selectedKeys={[pathname]}>
        <S.NavMenuItem key={GITHUB_URL}>
          <a href={GITHUB_URL} target="_blank" data-cy="github-link">
            GitHub
          </a>
        </S.NavMenuItem>
        <S.NavMenuItem key={DOCUMENTATION_URL}>
          <a href={DOCUMENTATION_URL} target="_blank" data-cy="documentation-link">
            Documentation
          </a>
        </S.NavMenuItem>
        <Menu.SubMenu
          key="help"
          icon={<QuestionCircleOutlined data-cy="onboarding-link" style={{color: '#61175e', fontSize: 16}} />}
        >
          <S.NavMenuItem key="guidedTour" onClick={handleGuidedTourCLick}>
            Show Onboarding
          </S.NavMenuItem>
        </Menu.SubMenu>
      </S.NavMenu>
    </S.Header>
  );
};

export default Header;
