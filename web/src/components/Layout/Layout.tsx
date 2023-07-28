import {ClusterOutlined, GlobalOutlined, SettingOutlined} from '@ant-design/icons';
import {Menu} from 'antd';
import React from 'react';
import {useLocation} from 'react-router-dom';

import logoAsset from 'assets/logo-white.svg';
import FileViewerModalProvider from 'components/FileViewerModal/FileViewerModal.provider';
import Header from 'components/Header';
import Link from 'components/Link';
import useRouterSync from 'hooks/useRouterSync';
import ConfirmationModalProvider from 'providers/ConfirmationModal';
import EnvironmentProvider from 'providers/Environment';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import MissingVariablesModalProvider from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import NotificationProvider from 'providers/Notification/Notification.provider';
import {ConfigMode} from 'types/DataStore.types';
import * as S from './Layout.styled';

interface IProps {
  children?: React.ReactNode;
  hasMenu?: boolean;
}

const menuItems = [
  {
    key: '0',
    icon: <ClusterOutlined />,
    label: <Link to="/">Tests</Link>,
    path: '/',
  },
  {
    key: '1',
    icon: <GlobalOutlined />,
    label: <Link to="/environments">Environments</Link>,
    path: '/environments',
  },
];

const footerMenuItems = [
  {
    key: '0',
    icon: <SettingOutlined />,
    label: <Link to="/settings">Settings</Link>,
    path: '/settings',
  },
];

const Layout = ({children, hasMenu = false}: IProps) => {
  useRouterSync();
  const {dataStoreConfig, isLoading} = useSettingsValues();
  const pathname = useLocation().pathname;
  const isNoTracingMode = dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;

  return (
    <NotificationProvider>
      <MissingVariablesModalProvider>
        <FileViewerModalProvider>
          <ConfirmationModalProvider>
            <EnvironmentProvider>
              <S.Layout hasSider>
                {hasMenu && (
                  <S.Sider width={256}>
                    <S.LogoContainer>
                      <Link to="/">
                        <img alt="Tracetest logo" src={logoAsset} />
                      </Link>
                    </S.LogoContainer>

                    <S.SiderContent>
                      <S.MenuContainer>
                        <Menu
                          defaultSelectedKeys={[
                            menuItems.findIndex(value => value.path === pathname).toString() || '0',
                          ]}
                          items={menuItems}
                          mode="inline"
                          theme="dark"
                        />
                      </S.MenuContainer>

                      <S.MenuContainer>
                        <Menu
                          defaultSelectedKeys={[
                            footerMenuItems.findIndex(value => value.path === pathname).toString() || '0',
                          ]}
                          items={footerMenuItems}
                          mode="inline"
                          theme="dark"
                        />
                      </S.MenuContainer>
                    </S.SiderContent>
                  </S.Sider>
                )}

                <S.Layout>
                  <Header hasLogo={!hasMenu} isNoTracingMode={isNoTracingMode && !isLoading} />
                  <S.Content $hasMenu={hasMenu}>{children}</S.Content>
                </S.Layout>
              </S.Layout>
            </EnvironmentProvider>
          </ConfirmationModalProvider>
        </FileViewerModalProvider>
      </MissingVariablesModalProvider>
    </NotificationProvider>
  );
};

export default Layout;
