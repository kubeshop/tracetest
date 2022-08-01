import {Typography} from 'antd';
import styled from 'styled-components';

export const Step = styled.div`
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
`;

export const FormContainer = styled.div`
  height: 100%;
`;

export const Title = styled(Typography.Title).attrs({level: 2})`
  && {
    margin-bottom: 24px;
  }
`;
