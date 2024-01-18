import {Typography} from 'antd';
import styled from 'styled-components';

export const TopContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0;
`;

export const ButtonsContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 23px;
  padding: 16px 0;
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

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 0;
  }
`;
