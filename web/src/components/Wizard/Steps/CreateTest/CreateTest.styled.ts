import {Typography} from 'antd';
import styled from 'styled-components';

export const ButtonContainer = styled.div`
  margin-top: 44px;

  display: flex;
  justify-content: flex-end;
`;

export const TabText = styled(Typography.Text)`
  && {
    margin: 0;
    color: ${({theme}) => theme.color.textSecondary};
    font-size: ${({theme}) => theme.size.sm};
  }
`;
