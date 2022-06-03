import {QuestionCircleOutlined} from '@ant-design/icons';
import {useTour} from '@reactour/tour';
import {Popover, Typography} from 'antd';
import {useMemo, useState} from 'react';
import {useParams} from 'react-router-dom';
import {DOCUMENTATION_URL, GITHUB_URL} from '../../constants/Common.constants';
import * as S from './Header.styled';
import {ShowOnboardingContent} from './ShowOnboardingContent';

interface IProps {
  pathname: string;
  onGuidedTourClick: () => void;
}

export const HeaderMenu = ({pathname, onGuidedTourClick}: IProps) => {
  const {setIsOpen} = useTour();
  const [visible, setVisible] = useState(false);
  const handleGuidedTourCLick = () => setVisible(o => !o);
  const content = useMemo(
    () => ShowOnboardingContent(onGuidedTourClick, setIsOpen, setVisible),
    [setIsOpen, setVisible, onGuidedTourClick]
  );
  const params = useParams();
  return (
    <Popover
      visible={visible}
      content={content}
      title={() => <Typography.Title level={5}>Take a quick tour of Tracetest?</Typography.Title>}
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
            label: (
              <QuestionCircleOutlined
                data-cy="onboarding-link"
                style={{color: '#61175e', fontSize: 16, opacity: !params.runId ? 0.5 : 1}}
              />
            ),
            disabled: !params.runId,
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
    </Popover>
  );
};
