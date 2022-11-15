import {MoreOutlined, SettingOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const ActionButton = styled(MoreOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;

export const ActionButtonRunView = styled(SettingOutlined)`
  color: ${({theme}) => theme.color.primary};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;
