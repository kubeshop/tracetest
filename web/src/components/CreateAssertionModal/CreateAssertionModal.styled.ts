import {Form} from 'antd';
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
