import {Tag} from 'antd';
import styled from 'styled-components';

export const SelectedTag = styled(Tag)<{$isLast: boolean}>`
  margin-right: ${({$isLast}) => ($isLast ? '8px' : '2px')};
  background: rgba(3, 24, 73, 0.1);
  border: none;
`;
