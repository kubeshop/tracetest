import {BulbOutlined, GithubOutlined, ReadOutlined} from '@ant-design/icons';
import {Dropdown, Menu, MenuProps, Space} from 'antd';
import {useMemo} from 'react';
import {useLocation} from 'react-router-dom';

import slackIcon from 'assets/slack.svg';
import {COMMUNITY_SLACK_URL, DOCUMENTATION_URL, GITHUB_URL} from 'constants/Common.constants';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import GuidedTourService from 'services/GuidedTour.service';
import Env from 'utils/Env';
import * as S from './Header.styled';

const appVersion = Env.get('appVersion');

function getMenuItems({
  isOnboardingActive,
  onClickOnboarding,
}: {
  isOnboardingActive: boolean;
  onClickOnboarding: () => void;
}) {
  const items: MenuProps['items'] = [
    {
      key: '1',
      label: (
        <a target="_blank" href={GITHUB_URL}>
          GitHub
        </a>
      ),
      icon: <GithubOutlined />,
    },
    {
      key: '2',
      label: (
        <a target="_blank" href={DOCUMENTATION_URL}>
          Documentation
        </a>
      ),
      icon: <ReadOutlined />,
    },
    {
      key: '3',
      label: (
        <a target="_blank" href={COMMUNITY_SLACK_URL}>
          Join our Slack community
        </a>
      ),
      icon: <img src={slackIcon} width={12} />,
    },
    {
      key: '4',
      label: (
        <a data-cy="menu-onboarding" onClick={onClickOnboarding}>
          Show Onboarding
        </a>
      ),
      icon: <BulbOutlined />,
      disabled: !isOnboardingActive,
    },
    {
      type: 'divider',
    },
    {
      key: '5',
      label: <S.AppVersion>App version: {appVersion}</S.AppVersion>,
      disabled: true,
    },
  ];

  return items;
}

const HelpMenu = () => {
  const pathname = useLocation().pathname;
  const tourByPathname = GuidedTourService.getByPathName(pathname);
  const isOnboardingActive = !!tourByPathname;
  const {onStart} = useGuidedTour();
  const items = useMemo(
    () => getMenuItems({isOnboardingActive, onClickOnboarding: onStart}),
    [isOnboardingActive, onStart]
  );

  return (
    <Dropdown overlay={<Menu items={items} />}>
      <a onClick={e => e.preventDefault()}>
        <Space>
          <S.QuestionIcon data-cy="menu-link" />
        </Space>
      </a>
    </Dropdown>
  );
};

export default HelpMenu;
