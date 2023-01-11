import styled from 'styled-components';
import {Typography} from 'antd';

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.md};
    margin-bottom: 16px;
    font-weight: 700;
  }
`;

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  margin-bottom: 24px;
  background: ${({theme}) => theme.color.white};
`;

export const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;
