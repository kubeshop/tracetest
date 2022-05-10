import styled from 'styled-components';
import Layout from 'antd/lib/layout';
import Title from 'antd/lib/typography/Title';
import Menu from 'antd/lib/menu';

export const Header = styled(Layout.Header)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  background: #fff;
  border-bottom: 1px solid #e2e4e6;
`;

export const TitleText = styled(Title).attrs({
  level: 2,
})`
  && {
    font-weight: 700;
    color: white;
    margin: 0;
  }
`;

export const NavMenu = styled(Menu).attrs({
  mode: 'horizontal',
  disabledOverflow: true,
})`
  && {
    align-items: center;
  }

  .ant-menu-item > span > a {
    color: #61175e;
  }
`;

export const NavMenuItem = styled(Menu.Item)``;
