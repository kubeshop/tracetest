import {Tag} from 'antd';
import styled from 'styled-components';

export const Badge = styled(Tag).attrs(({theme}) => ({
  color: theme.color.borderLight,
}))`
  font-size: ${({theme}) => theme.size.md};
  margin: 0;
  margin-left: 8px;
  padding: 0 8px;
  line-height: 15px;
  border-radius: 8px;
  font-weight: 400;
  background: ${({theme}) => theme.color.textSecondary};
  color: ${({theme}) => theme.color.primary};
`;
