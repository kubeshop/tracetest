import styled from 'styled-components';
import {Input} from 'antd';

export const SearchInput = styled(Input)<{height: string; width: string}>`
  height: ${({height}) => height};
  width: ${({width}) => width};
`;
