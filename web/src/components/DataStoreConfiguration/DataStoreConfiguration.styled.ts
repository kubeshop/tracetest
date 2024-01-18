import {Typography} from 'antd';
import styled from 'styled-components';

export const TopContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 22px;
`;

export const ButtonsContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 23px;
  padding: 16px 22px;
  border-top: 1px solid ${({theme}) => theme.color.borderLight};
  position: sticky;
  bottom: 0;
  background: white;
`;

export const Description = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
  }
`;