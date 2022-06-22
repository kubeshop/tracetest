import {Form, FormInstance} from 'antd';
import React from 'react';
import {IRequestDetailsValues} from '../RequestDetails';
import {apiKeyFields} from './apiKeyFields';
import {basicFields} from './basicFields';
import {bearerFields} from './bearerFields';
import {CreateTestFormAuthTypeInput} from './CreateTestFormAuthTypeInput';

const RequestDetailsAuthInput: React.FC<{form: FormInstance<IRequestDetailsValues>}> = ({form}) => {
  return (
    <div>
      <CreateTestFormAuthTypeInput form={form} />
      <Form.Item noStyle style={{width: '100%'}} shouldUpdate>
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
