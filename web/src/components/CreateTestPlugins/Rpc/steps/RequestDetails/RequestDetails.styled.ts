import {Button} from 'antd';
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

export const InputContainer = styled.div`
  display: grid;
  gap: 26px;
  grid-template-columns: 50%;
  justify-content: start;
`;

export const DoubleInputContainer = styled(InputContainer)`
  display: grid;
  gap: 26px;
  grid-template-columns: 50% 50%;
`;

export const Row = styled.div`
  display: flex;
`;

export const URLInputContainer = styled.div`
  display: flex;
  align-items: end;
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  margin-bottom: 8px;
`;

export const UploadButton = styled(Button).attrs({
  type: 'primary',
  ghost: true,
})`
  width: 490px;
`;
