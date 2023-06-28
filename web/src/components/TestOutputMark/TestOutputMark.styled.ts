import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  background-color: ${({theme}) => theme.color.warningYellow};
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  height: 15px;
  justify-content: center;
  width: 15px;
`;

export const Text = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.white};
    font-size: ${({theme}) => theme.size.sm};
    font-weight: bold;
    line-height: 1;
  }
`;
