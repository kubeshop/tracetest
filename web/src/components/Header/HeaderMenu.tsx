import {Popover, Typography} from 'antd';

import {DOCUMENTATION_URL, GITHUB_URL} from 'constants/Common.constants';
import {useLocation} from 'react-router-dom';
import * as S from './Header.styled';

const HeaderMenu = () => {
  const {pathname} = useLocation();

  return (
    <Popover
      arrowContent={null}
      content={() => null}
      title={() => <Typography.Title level={2}>Take a quick tour of Tracetest?</Typography.Title>}
      visible={false}
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
                disabled: true,
                label: (
                  <a key="guidedTour" aria-disabled>
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
