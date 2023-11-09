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

export const BodyFieldContainer = styled.div<{$isDisplaying: boolean}>`
  width: 100%;
  display: ${({$isDisplaying}) => ($isDisplaying ? 'none' : 'unset')};

  && {
    .cm-editor {
      overflow: hidden;
      display: flex;
      border-radius: 2px;
      font-size: ${({theme}) => theme.size.md};
      outline: 1px solid grey;
      font-family: SFPro, Inter, serif;
    }

    .cm-line {
      padding: 0;

      span {
        font-family: SFPro, Inter, serif;
      }
    }
  }
`;
