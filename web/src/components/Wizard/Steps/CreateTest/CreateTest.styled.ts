import {Typography} from 'antd';
import styled from 'styled-components';

export const ButtonContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  margin-top: 4px;
`;

export const TabText = styled(Typography.Text)`
  && {
    margin: 0;
    color: ${({theme}) => theme.color.textSecondary};
    font-size: ${({theme}) => theme.size.sm};
  }
`;
