import {PlusOutlined} from '@ant-design/icons';
import {Button, Form, Input} from 'antd';
import React from 'react';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import * as S from './RequestDetails.styled';

const RequestDetailsHeadersInput: React.FC = () => (
  <Form.Item className="input-headers" label="Headers list">
    <Form.List name="headers" initialValue={DEFAULT_HEADERS}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.HeaderContainer key={field.name}>
              <Form.Item name={[field.name, 'key']} noStyle>
                <Input placeholder={`Header ${index + 1}`} />
              </Form.Item>

              <Form.Item name={[field.name, 'value']} noStyle>
                <Input placeholder={`Value ${index + 1}`} />
              </Form.Item>

              <Form.Item noStyle>
                <Button
                  icon={<S.DeleteIcon />}
                  onClick={() => remove(field.name)}
                  style={{marginLeft: 12}}
                  type="link"
                />
              </Form.Item>
            </S.HeaderContainer>
          ))}

          <Button
            data-cy="add-header"
            icon={<PlusOutlined />}
            onClick={() => add()}
            style={{fontWeight: 600, height: 'auto', padding: 0}}
            type="link"
          >
            Add Header
          </Button>
        </>
      )}
    </Form.List>
  </Form.Item>
);

export default RequestDetailsHeadersInput;
