import {useTour} from '@reactour/tour';
import {Popover, Typography} from 'antd';
import {useMemo} from 'react';
import {useParams} from 'react-router-dom';

import {DOCUMENTATION_URL, GITHUB_URL} from 'constants/Common.constants';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import * as S from './Header.styled';
import {ShowOnboardingContent} from './ShowOnboardingContent';

interface IProps {
  pathname: string;
  onGuidedTourClick: () => void;
}

export const HeaderMenu = ({pathname, onGuidedTourClick}: IProps) => {
  const {setIsOpen} = useTour();
  const {setIsTriggerVisible, isTriggerVisible} = useGuidedTour();

  const content = useMemo(
    () =>
      ShowOnboardingContent(
        onGuidedTourClick,
        () => setIsOpen(true),
        () => setIsTriggerVisible(false)
      ),
    [onGuidedTourClick, setIsOpen, setIsTriggerVisible]
  );

  const params = useParams();
  return (
    <Popover
      visible={isTriggerVisible}
      content={content}
      title={() => <Typography.Title level={2}>Take a quick tour of Tracetest?</Typography.Title>}
      arrowContent={null}
    >
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
            label: <S.QuestionIcon $disabled={!params.runId} data-cy="onboarding-link" />,
            disabled: !params.runId,
            children: [
              {
                key: 'Onboarding',
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
