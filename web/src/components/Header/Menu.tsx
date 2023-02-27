import {useLocation} from 'react-router-dom';

import {DISCORD_URL, DOCUMENTATION_URL, GITHUB_URL} from 'constants/Common.constants';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import Env from 'utils/Env';
import * as S from './Header.styled';
import GuidedTourService from '../../services/GuidedTour.service';

const appVersion = Env.get('appVersion');

const Menu = () => {
  const pathname = useLocation().pathname;
  const tourByPathname = GuidedTourService.getByPathName(pathname);
  const {onStart} = useGuidedTour();

  return (
    <S.NavMenu
      items={[
        {
          key: 'SubMenu',
          label: <S.QuestionIcon data-cy="menu-link" />,
          children: [
            {
              key: 'github',
              label: (
                <a data-cy="github-link" href={GITHUB_URL} target="_blank">
                  GitHub
                </a>
              ),
            },
            {
              key: 'docs',
              label: (
                <a data-cy="documentation-link" href={DOCUMENTATION_URL} target="_blank">
                  Documentation
                </a>
              ),
            },
            {
              key: 'discord',
              label: (
                <a data-cy="discord-link" href={DISCORD_URL} target="_blank">
                  Discord
                </a>
              ),
            },
            {
              key: 'Onboarding',
              disabled: !tourByPathname,
              label: (
                <a key="guidedTour" onClick={onStart}>
                  Show Onboarding
                </a>
              ),
            },
            {
              key: 'App version',
              disabled: true,
              label: (
                <S.AppVersionContainer>
                  <S.AppVersion>App version: {appVersion}</S.AppVersion>
                </S.AppVersionContainer>
              ),
            },
          ],
        },
      ]}
    />
  );
};

export default Menu;
