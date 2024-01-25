import styled from 'styled-components';
import {Typography} from 'antd';

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.md};
    margin-bottom: 16px;
    font-weight: 600;
  }
`;

export const Wrapper = styled.div`
  background: ${({theme}) => theme.color.white};
  border-radius: 2px;
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  display: flex;
  flex-direction: column;
`;

export const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;
