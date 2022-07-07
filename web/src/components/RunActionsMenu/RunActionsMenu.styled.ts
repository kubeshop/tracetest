import {MoreOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const ActionButton = styled(MoreOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  font-size: 24px;
`;
