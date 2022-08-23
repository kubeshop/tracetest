import {SettingOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const ActionButton = styled(SettingOutlined)`
  color: ${({theme}) => theme.color.primary};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;
