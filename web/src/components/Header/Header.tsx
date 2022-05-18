import {QuestionCircleOutlined} from '@ant-design/icons';
import {useTour} from '@reactour/tour';
import Logo from 'assets/Logo.svg';
import {DOCUMENTATION_URL, GITHUB_URL} from 'constants/Common.constants';
import {FC} from 'react';
import {Link, useLocation} from 'react-router-dom';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
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
          <img alt="tracetest_log" data-cy="logo" src={Logo} />
        </S.TitleText>
      </Link>
      <S.NavMenu
        selectedKeys={[pathname]}
        items={[
          {
            key: 'github',
            label: (
              <a href={GITHUB_URL} target="_blank" data-cy="github-link">
                GitHub
              </a>
            ),
          },
          {
            key: 'docs',
            label: (
              <a href={DOCUMENTATION_URL} target="_blank" data-cy="documentation-link">
                Documentation
              </a>
            ),
          },
          {
            key: 'SubMenu',
            label: <QuestionCircleOutlined data-cy="onboarding-link" style={{color: '#61175e', fontSize: 16}} />,
            children: [
              {
                key: 'Onboarding',
                label: (
                  <a key="guidedTour" onClick={handleGuidedTourCLick}>
                    Show Onboarding
                  </a>
                ),
              },
            ],
          },
        ]}
      />
    </S.Header>
  );
};

export default Header;
