import {QuestionCircleOutlined} from '@ant-design/icons';
import {Layout, Menu} from 'antd';
import styled from 'styled-components';

export const Header = styled(Layout.Header)`
  align-items: center;
  background: ${({theme}) => theme.color.white};
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  display: flex;
  justify-content: space-between;
  height: 48px;
  line-height: 48px;
  padding: 0 24px;

  .ant-dropdown-trigger {
    display: block;
  }
  .ant-popover-title {
    padding: 16px;
  }
  .ant-popover-arrow {
    display: none;
  }
  .ant-popover-inner-content {
    padding: 0;
  }
`;

export const Logo = styled.img`
  height: 24px;
  margin-right: 24px;
`;

export const NavMenu = styled(Menu).attrs({
  mode: 'horizontal',
  disabledOverflow: true,
})`
  && {
    align-items: center;
  }

  .ant-menu-submenu.ant-menu-submenu-horizontal {
    padding-right: 0;
  }

  .ant-menu-item > span > a {
    color: ${({theme}) => theme.color.primary};
  }
`;

export const QuestionIcon = styled(QuestionCircleOutlined)`
  color: ${({theme}) => theme.color.primary};
  font-size: ${({theme}) => theme.size.lg};
`;
