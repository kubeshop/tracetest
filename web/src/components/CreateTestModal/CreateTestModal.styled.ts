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
