import {Typography} from 'antd';
import styled from 'styled-components';

export const Row = styled.div`
  display: flex;
`;

export const FormSection = styled.div`
  background-color: ${({theme}) => theme.color.white};
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 4px;
  margin-bottom: 16px;
  padding: 16px;
`;

export const FormSectionTitle = styled(Typography.Text).attrs({strong: true})`
  margin-bottom: 8px;
  display: block;
`;
