import styled from 'styled-components';
import {Button, Typography} from 'antd';

export const Container = styled.div`
  background-color: ${({theme}) => theme.color.background};
  border-radius: 8px;
  left: 16px;
  padding: 0 8px;
  position: absolute;
  top: 16px;
  z-index: 9;
`;

export const NavigationText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.primary};
  margin-left: 8px;
`;

export const ToggleButton = styled(Button)<{$isActive?: boolean}>`
  color: ${({theme}) => theme.color.primary};
  width: 24px;

  &:focus {
    background-color: unset;
    color: ${({theme}) => theme.color.primary};
  }

  &:hover {
    background-color: unset;
    color: ${({theme}) => theme.color.primary};
  }
`;
