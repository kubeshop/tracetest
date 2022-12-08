import styled from 'styled-components';
import {Typography} from 'antd';

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.md};
    margin-bottom: 16px;
    font-weight: 700;
  }
`;

export const Description = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
  }
`;

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  margin-top: 24px;
  background: ${({theme}) => theme.color.white};
`;

export const FormContainer = styled.div`
  padding: 16px;
  min-height: 663px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
`;

export const ButtonsContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 23px;
  padding-top: 16px;
  border-top: 1px solid ${({theme}) => theme.color.borderLight};
`;
