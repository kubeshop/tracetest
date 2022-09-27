import {PlusOutlined} from '@ant-design/icons';
import {Button, Form, Input} from 'antd';
import {DEFAULT_HEADERS, Header} from 'constants/Test.constants';
import React from 'react';
import * as S from './RequestDetails.styled';

interface IProps {
  initialValue?: Header[];
  name?: string;
  unit?: string;
  addLabel?: string;
}
const RequestDetailsHeadersInput: React.FC<IProps> = ({
  unit = 'Header',
  name = 'headers',
  initialValue = DEFAULT_HEADERS,
  addLabel = 'Add Header',
}) => (
  <Form.Item
    className="input-headers"
    label={`${name?.replace(/(^\w{1})|(\s+\w{1})/g, letter => letter.toUpperCase())} list`}
    shouldUpdate
  >
    <Form.List name={name} initialValue={initialValue}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.HeaderContainer key={field.name}>
              <Form.Item name={[field.name, 'key']} noStyle>
                <Input placeholder={`${unit} ${index + 1}`} />
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
            {addLabel}
          </Button>
        </>
      )}
    </Form.List>
  </Form.Item>
);

export default RequestDetailsHeadersInput;
