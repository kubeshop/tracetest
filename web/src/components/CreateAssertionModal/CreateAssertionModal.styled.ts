import {Form, Select as AntSelect} from 'antd';
import styled from 'styled-components';

export const FullHeightFormItem = styled(Form.Item)`
  .ant-select-single:not(.ant-select-customize-input) .ant-select-selector {
    height: 40px;
    display: flex;
    align-items: center;
  }

  input {
    height: 40px;
  }
`;

export const Select = styled(AntSelect)`
  min-width: 88px;
  > .ant-select-selector {
    min-height: 100%;
  }
`;
