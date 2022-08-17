import {useTour} from '@reactour/tour';
import {Popover, Typography} from 'antd';
import {useMemo} from 'react';
import {useLocation, useParams} from 'react-router-dom';

import {DOCUMENTATION_URL, GITHUB_URL} from 'constants/Common.constants';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import * as S from './Header.styled';
import {ShowOnboardingContent} from './ShowOnboardingContent';

const {onGuidedTourClick} = HomeAnalyticsService;

export const HeaderMenu = () => {
  const {pathname} = useLocation();
  const params = useParams();
  const {setIsOpen} = useTour();
  const {setIsTriggerVisible, isTriggerVisible} = useGuidedTour();

  const content = useMemo(
    () =>
      ShowOnboardingContent(
        onGuidedTourClick,
        () => setIsOpen(true),
        () => setIsTriggerVisible(false)
      ),
    [setIsOpen, setIsTriggerVisible]
  );

  return (
    <Popover
      arrowContent={null}
      content={content}
      title={() => <Typography.Title level={2}>Take a quick tour of Tracetest?</Typography.Title>}
      visible={isTriggerVisible}
    >
      <S.NavMenu
        selectedKeys={[pathname]}
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
                key: 'Onboarding',
                disabled: !params.runId,
                label: (
                  <a key="guidedTour" onClick={() => setIsTriggerVisible(!isTriggerVisible)}>
                    Show Onboarding
                  </a>
                ),
              },
            ],
          },
        ]}
      />
    </Popover>
  );
};
