import styled from 'styled-components';
import {Button as AntdButton} from 'antd';

export const DemoContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

export const Button = styled(AntdButton).attrs({
  size: 'small',
})`
  && {
    border-radius: 16px;
    border: none;
    background: ${({theme}) => theme.color.primary};
    color: ${({theme}) => theme.color.white};
  }
`;
