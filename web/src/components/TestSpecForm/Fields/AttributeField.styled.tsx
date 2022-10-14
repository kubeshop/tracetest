import styled from 'styled-components';
import {Form, Tag as AntdTag, Select as AntSelect} from 'antd';

export const Container = styled.div`
  display: inline-flex;
  width: 100%;
  .ant-select-item-option-active:hover {
    background-color: #faf4fe !important;
  }
`;

export const Select = styled(AntSelect)`
  min-width: 88px;

  > .ant-select-selector {
    min-height: 100%;
  }
`;

export const FormItem = styled(Form.Item)`
  margin: 0;
`;
export const Sides = styled.div`
  flex-basis: 50%;
  max-width: 50%;
  min-width: 50%;
`;

export const Tag = styled(AntdTag)`
  background: #e7e8eb;
`;

export const Title = styled.p`
  font-weight: bold;
`;

export const TagsContainer = styled.div`
  margin-top: 8px;
`;

export const Content = styled.div`
  padding: 10px 10px;
`;
