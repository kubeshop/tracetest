import {Typography} from 'antd';
import styled, {createGlobalStyle} from 'styled-components';

export const GlobalStyle = createGlobalStyle`
  .select-method.ant-select .ant-select-selector,
  .select-dropdown-method .ant-select-item-option-selected {
    background-color: #fafafa;
  }

  .select-method .ant-select-arrow,
  .input-name .ant-form-item-label > label,
  .input-headers .ant-form-item-label > label,
  .input-body .ant-form-item-label > label {
    color: #031849;
  }
`;

export const DemoContainer = styled.div`
  margin-bottom: 24px;
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  margin-bottom: 8px;
`;

export const Row = styled.div`
  display: flex;
`;

export const FormSection = styled.div`
  border: 1px solid #e8e8e8;
  padding: 16px;
  border-radius: 4px;
  margin-bottom: 16px;
`;

export const FormSectionTitle = styled(Typography.Text).attrs({strong: true})`
  margin-bottom: 8px;
  display: block;
`;
