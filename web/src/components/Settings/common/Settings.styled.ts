import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  background: ${({theme}) => theme.color.white};
  display: flex;
  flex-direction: column;
  margin-bottom: 24px;
  padding: 16px;
`;

export const FormContainer = styled.div`
  margin: 16px 0;
`;

export const FooterContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
  padding: 16px 22px 0;
  border-top: 1px solid ${({theme}) => theme.color.borderLight};
`;

export const Description = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
  }
`;
