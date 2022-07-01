import {Form, FormInstance} from 'antd';
import React from 'react';
import {TRequestAuth} from 'types/Test.types';
import {apiKeyFields} from './apiKeyFields';
import {basicFields} from './basicFields';
import {bearerFields} from './bearerFields';
import TypeInput from './TypeInput';

interface IProps {
  form: FormInstance<{auth: TRequestAuth}>;
}

const RequestDetailsAuthInput = ({form}: IProps) => {
  return (
    <div>
      <TypeInput form={form} />
      <Form.Item noStyle shouldUpdate style={{marginBottom: 0, width: '100%'}}>
        {({getFieldValue}) => {
          const method = getFieldValue('auth')?.type;
          switch (method) {
            case 'bearer':
              return bearerFields;
            case 'basic':
              return basicFields;
            case 'apiKey':
              return apiKeyFields;
            default:
              return null;
          }
        }}
      </Form.Item>
    </div>
  );
};

export default RequestDetailsAuthInput;
