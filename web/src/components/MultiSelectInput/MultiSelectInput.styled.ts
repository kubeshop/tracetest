import {Tag} from 'antd';
import styled from 'styled-components';

export const SelectedTag = styled(Tag)<{$isLast: boolean}>`
  background: ${({theme}) => theme.color.borderLight};
  border: none;
  margin-right: ${({$isLast}) => ($isLast ? '8px' : '2px')};
`;
