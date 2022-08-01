import {SearchOutlined} from '@ant-design/icons';
import {Input} from 'antd';
import styled from 'styled-components';

export const SearchInput = styled(Input)<{height: string; width: string}>`
  height: ${({height}) => height};
  width: ${({width}) => width};
`;

export const SearchIcon = styled(SearchOutlined)`
  color: ${({theme}) => theme.color.textLight};
`;
