import {ClusterOutlined, GlobalOutlined} from '@ant-design/icons';
import {Menu} from 'antd';
import {MenuInfo} from 'rc-menu/es/interface';
import React from 'react';
import {Link, useNavigate} from 'react-router-dom';

import logoAsset from 'assets/logo-white.svg';
import FileViewerModalProvider from 'components/FileViewerModal/FileViewerModal.provider';
import Header from 'components/Header';
import useRouterSync from 'hooks/useRouterSync';
import ConfirmationModalProvider from 'providers/ConfirmationModal';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import * as S from './Layout.styled';

interface IProps {
  children?: React.ReactNode;
  hasMenu?: boolean;
}

const menuItems = [
  {
    key: '0',
    icon: <ClusterOutlined />,
    label: 'Tests',
    path: '/',
  },
];

if (ExperimentalFeature.isEnabled('transactions')) {
  menuItems.push({
    key: '1',
    icon: <GlobalOutlined />,
    label: 'Environments',
    path: '/environments',
  });
}

const Layout = ({children, hasMenu = false}: IProps) => {
  const navigate = useNavigate();
  useRouterSync();

  const handleOnClickMenu = (menuInfo: MenuInfo) => {
    const item = menuItems.find(menuItem => menuItem.key === menuInfo.key);
    navigate(item?.path ?? '/');
  };

  return (
    <FileViewerModalProvider>
      <ConfirmationModalProvider>
        <S.Layout hasSider>
          {hasMenu && (
            <S.Sider width={256}>
              <S.LogoContainer>
                <Link to="/">
                  <img alt="Tracetest logo" src={logoAsset} />
                </Link>
              </S.LogoContainer>

              <S.MenuContainer>
                <Menu
                  defaultSelectedKeys={['0']}
                  items={menuItems}
                  mode="inline"
                  onClick={handleOnClickMenu}
                  theme="dark"
                />
              </S.MenuContainer>
            </S.Sider>
          )}

          <S.Layout>
            <Header hasEnvironments={ExperimentalFeature.isEnabled('transactions')} hasLogo={!hasMenu} />
            <S.Content $hasMenu={hasMenu}>{children}</S.Content>
          </S.Layout>
        </S.Layout>
      </ConfirmationModalProvider>
    </FileViewerModalProvider>
  );
};

export default Layout;
