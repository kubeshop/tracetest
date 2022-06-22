import {Typography} from 'antd';
import styled from 'styled-components';

export const Step = styled.div`
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
`;

export const FormContainer = styled.div`
  height: 100%;
`;

export const Title = styled(Typography.Title).attrs({
  level: 5,
})`
  && {
    margin-bottom: 24px;
  }
`;
