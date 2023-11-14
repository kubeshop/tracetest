import {Typography} from 'antd';
import styled from 'styled-components';

export const Row = styled.div`
  display: flex;
`;

export const Label = styled(Typography.Text).attrs({as: 'div'})`
  margin-bottom: 8px;
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: grid;
  justify-content: center;
  grid-template-columns: 40% 40% 19%;
  margin-bottom: 8px;
`;

export const SSLVerificationContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
`;
