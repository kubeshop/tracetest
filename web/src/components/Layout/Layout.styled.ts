import {Layout as LayoutAntd} from 'antd';
import styled, {css} from 'styled-components';

export const Content = styled(LayoutAntd.Content)<{$hasMenu: boolean}>`
  display: flex;
  flex-direction: column;

  ${({$hasMenu}) =>
    $hasMenu &&
    css`
      height: 100%;
      overflow-y: scroll;
    `}
`;

export const Layout = styled(LayoutAntd)`
  min-height: 100%;
  background: ${({theme}) => theme.color.background};
`;

export const LogoContainer = styled.div`
  background: #2c1a54;
  padding: 7px 55px;
`;

export const MenuContainer = styled.div`
  .ant-menu.ant-menu-dark {
    background: transparent;
  }

  .ant-menu-dark.ant-menu-dark:not(.ant-menu-horizontal) .ant-menu-item-selected {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
  }

  padding: 24px 20px;
`;

export const Sider = styled(LayoutAntd.Sider)`
  .ant-layout-sider-children {
    display: flex;
    flex-direction: column;
  }

  .ant-layout-sider-trigger {
    background: #2c1a54;
  }

  background: linear-gradient(180deg, #2f1e61 0%, #8b2c53 111.31%, rgba(49, 38, 132, 0) 180.18%, #df4f80 180.18%);
`;

export const SiderContent = styled.div`
  display: flex;
  flex: 1;
  flex-direction: column;
  justify-content: space-between;
`;
