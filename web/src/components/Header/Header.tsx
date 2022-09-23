import {DownOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Space} from 'antd';

import {Link} from 'react-router-dom';

import Logo from 'assets/Logo.svg';
import * as S from './Header.styled';
import HeaderMenu from './HeaderMenu';

interface IProps {
  hasEnvironments?: boolean;
  hasLogo?: boolean;
}

const menu = (
  <Menu
    items={[
      {
        key: '1',
        label: (
          <a target="_blank" rel="noopener noreferrer" href="#">
            Env 1
          </a>
        ),
      },
      {
        key: '2',
        label: (
          <a target="_blank" rel="noopener noreferrer" href="#">
            Env 2
          </a>
        ),
      },
    ]}
  />
);

const Header = ({hasEnvironments = false, hasLogo = false}: IProps) => (
  <S.Header>
    <div>
      {hasLogo && (
        <Link to="/">
          <S.Logo alt="tracetest logo" data-cy="logo" src={Logo} />
        </Link>
      )}
    </div>

    <Space>
      {hasEnvironments && (
        <Dropdown overlay={menu}>
          <a onClick={e => e.preventDefault()}>
            <Space>
              All environments
              <DownOutlined />
            </Space>
          </a>
        </Dropdown>
      )}

      <HeaderMenu />
    </Space>
  </S.Header>
);

export default Header;
