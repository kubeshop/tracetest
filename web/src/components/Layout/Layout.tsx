import {AppstoreAddOutlined, ClusterOutlined, GlobalOutlined, HomeOutlined, SettingOutlined} from '@ant-design/icons';
import {Menu} from 'antd';
import {Outlet, useLocation} from 'react-router-dom';

import logoAsset from 'assets/logo-white.svg';
import FileViewerModalProvider from 'components/FileViewerModal/FileViewerModal.provider';
import Header from 'components/Header';
import Link from 'components/Link';
import useRouterSync from 'hooks/useRouterSync';
import ConfirmationModalProvider from 'providers/ConfirmationModal';
import CreateTestProvider from 'providers/CreateTest';
import GuidedTourProvider from 'providers/GuidedTour';
import VariableSetProvider from 'providers/VariableSet';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import MissingVariablesModalProvider from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import NotificationProvider from 'providers/Notification/Notification.provider';
import {ConfigMode} from 'types/DataStore.types';
import * as S from './Layout.styled';
import MenuBottom from './MenuBottom';

export type TCustomHeader = typeof Header;

interface IProps {
  hasMenu?: boolean;
}

const menuItems = [
  {
    key: '0',
    icon: <HomeOutlined />,
    label: <Link to="/">Home</Link>,
    path: '/',
  },
  {
    key: '1',
    icon: <ClusterOutlined />,
    label: <Link to="/tests">Tests</Link>,
    path: '/tests',
  },
  {
    key: '2',
    icon: <AppstoreAddOutlined />,
    label: <Link to="/testsuites">Test Suites</Link>,
    path: '/testsuites',
  },
  {
    key: '3',
    icon: <GlobalOutlined />,
    label: <Link to="/variablesets">Variable Sets</Link>,
    path: '/variablesets',
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

const Layout = ({hasMenu = false}: IProps) => {
  useRouterSync();
  const {dataStoreConfig, isLoading} = useSettingsValues();
  const pathname = useLocation().pathname;
  const isNoTracingMode = dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;

  return (
    <ConfirmationModalProvider>
      <NotificationProvider>
        <MissingVariablesModalProvider>
          <FileViewerModalProvider>
            <VariableSetProvider>
              <GuidedTourProvider>
                <CreateTestProvider>
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
                            <MenuBottom />
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
                      <S.Content $hasMenu={hasMenu}>
                        <Outlet />
                      </S.Content>
                    </S.Layout>
                  </S.Layout>
                </CreateTestProvider>
              </GuidedTourProvider>
            </VariableSetProvider>
          </FileViewerModalProvider>
        </MissingVariablesModalProvider>
      </NotificationProvider>
    </ConfirmationModalProvider>
  );
};

export default Layout;
