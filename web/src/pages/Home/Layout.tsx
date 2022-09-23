import {ClusterOutlined, GlobalOutlined} from '@ant-design/icons';
import {Layout as LayoutAntd, Menu} from 'antd';
import {MenuInfo} from 'rc-menu/es/interface';
import React from 'react';
import {Link, useNavigate} from 'react-router-dom';

import logoAsset from 'assets/logo-white.svg';
import Header from 'components/Header';
import * as S from './Layout.styled';

interface IProps {
  children?: React.ReactNode;
}

const menuItems = [
  {
    key: '0',
    icon: <ClusterOutlined />,
    label: 'Tests',
    path: '/',
  },
  {
    key: '1',
    icon: <GlobalOutlined />,
    label: 'Environments',
    path: '/environments',
  },
];

const Layout = ({children}: IProps) => {
  const navigate = useNavigate();

  const handleOnClickMenu = (menuInfo: MenuInfo) => {
    const item = menuItems.find(menuItem => menuItem.key === menuInfo.key);
    navigate(item?.path ?? '/');
  };

  return (
    <LayoutAntd hasSider>
      <S.Sider width={256}>
        <S.LogoContainer>
          <Link to="/">
            <img alt="Tracetest logo" src={logoAsset} />
          </Link>
        </S.LogoContainer>

        <S.MenuContainer>
          <Menu defaultSelectedKeys={['0']} items={menuItems} mode="inline" onClick={handleOnClickMenu} theme="dark" />
        </S.MenuContainer>
      </S.Sider>

      <LayoutAntd>
        <Header hasEnvironments />
        <S.Content>{children}</S.Content>
      </LayoutAntd>
    </LayoutAntd>
  );
};

export default Layout;
