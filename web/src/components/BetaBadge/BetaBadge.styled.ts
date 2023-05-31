import {Tag} from 'antd';
import styled from 'styled-components';

export const Badge = styled(Tag).attrs(({theme}) => ({
  color: theme.color.primary,
}))`
  font-size: 10px;
  margin-left: 5px;
  padding: 0 2px;
  margin-right: 0;
  line-height: 15px;
`;
