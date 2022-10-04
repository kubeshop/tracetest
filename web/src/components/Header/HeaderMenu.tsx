import {Popover, Typography} from 'antd';

import {DOCUMENTATION_URL, GITHUB_URL} from 'constants/Common.constants';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import {useMemo} from 'react';
import {useLocation, useParams} from 'react-router-dom';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {switchTraceMode} from '../GuidedTour/traceStepList';
import * as S from './Header.styled';
import {ShowOnboardingContent} from './ShowOnboardingContent';

const {onGuidedTourClick} = HomeAnalyticsService;

const HeaderMenu = () => {
  const {pathname} = useLocation();
  const params = useParams();
  const {setState, state} = useGuidedTour();

  const content = useMemo(
    () =>
      ShowOnboardingContent(
        onGuidedTourClick,
        () => {
          switchTraceMode(0)();
          setState(st => ({...st, tourActive: true, run: true}));
        },
        () => setState(st => ({...st, dialog: false}))
      ),
    [setState]
  );

  return (
    <Popover
      arrowContent={null}
      content={content}
      title={() => <Typography.Title level={2}>Take a quick tour of Tracetest?</Typography.Title>}
      visible={state.dialog}
      trigger={['click']}
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
                  <a key="guidedTour" onClick={() => setState(st => ({...st, dialog: !st.dialog}))}>
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

export default HeaderMenu;
