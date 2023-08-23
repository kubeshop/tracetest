import {QuestionCircleOutlined} from '@ant-design/icons';
import {Layout, Typography} from 'antd';
import styled from 'styled-components';

export const Header = styled(Layout.Header)`
  align-items: center;
  background: ${({theme}) => theme.color.white};
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  display: flex;
  justify-content: space-between;
  height: 48px;
  line-height: 48px;
  padding: 0;
  padding-left: 24px;
  padding-right: 24px;

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
  margin-right: 12px;
`;

export const QuestionIcon = styled(QuestionCircleOutlined)`
  color: ${({theme}) => theme.color.primary};
  font-size: ${({theme}) => theme.size.lg};
`;

export const AppVersion = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.textLight};
    font-size: ${({theme}) => theme.size.sm};
  }
`;

export const MenuContainer = styled.div`
  display: flex;
  gap: 12px;
  flex-direction: row;
  align-items: center;
`;
