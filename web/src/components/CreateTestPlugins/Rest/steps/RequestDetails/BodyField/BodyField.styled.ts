import {Form} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  margin-bottom: 24px;
`;

export const RadioContainer = styled.div`
  margin-bottom: 8px;
`;

export const SingleLabel = styled(Form.Item)`
  && {
    margin-bottom: 0;

    .ant-form-item-control-input {
      display: none;
    }

    .ant-form-item-label {
      padding: 0 0 8px;
    }
  }
`;
