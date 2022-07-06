import {QuestionCircleOutlined} from '@ant-design/icons';
import {Layout, Menu} from 'antd';
import styled from 'styled-components';

export const Header = styled(Layout.Header)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  background: ${({theme}) => theme.color.white};
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};

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

export const NavMenu = styled(Menu).attrs({
  mode: 'horizontal',
  disabledOverflow: true,
})`
  && {
    align-items: center;
  }

  .ant-menu-item > span > a {
    color: ${({theme}) => theme.color.primary};
  }
`;

export const QuestionIcon = styled(QuestionCircleOutlined)<{$disabled: boolean}>`
  color: ${({theme}) => theme.color.primary};
  font-size: ${({theme}) => theme.size.lg};
  opacity: ${({$disabled}) => ($disabled ? 0.5 : 1)};
`;
