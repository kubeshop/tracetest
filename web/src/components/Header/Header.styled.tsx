import styled from 'styled-components';
import Layout from 'antd/lib/layout';
import Title from 'antd/lib/typography/Title';
import Menu from 'antd/lib/menu';

export const Header = styled(Layout.Header).attrs({
  theme: 'dark',
})`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
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
  theme: 'dark',
  mode: 'horizontal',
  disabledOverflow: true,
})`
  && {
    align-items: center;
  }

  .ant-menu-item > span > a {
    color: white;
  }
`;

export const NavMenuItem = styled(Menu.Item)``;
