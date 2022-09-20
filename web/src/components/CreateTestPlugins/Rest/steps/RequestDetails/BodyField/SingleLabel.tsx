import {Form} from 'antd';
import styled from 'styled-components';

export const SingleLabel = styled(Form.Item)`
  && {
    .ant-form-item-control-input {
      display: none;
    }

    .ant-form-item-label {
      padding: 0 0 8px;
    }
  }
`;
