import styled from 'styled-components';
import {Typography} from 'antd';

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin-bottom: 25px;
    font-weight: 700;
  }
`;

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
`;

export const FormContainer = styled.div`
  padding: 24px;
`;

export const ButtonsContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 23px;
`;
