import {Divider, Form, FormInstance} from 'antd';
import React from 'react';
import {ICreateTestValues} from '../CreateTestForm';
import {apiKeyFields} from './apiKeyFields';
import {basicFields} from './basicFields';
import {bearerFields} from './bearerFields';
import {CreateTestFormAuthTypeInput} from './CreateTestFormAuthTypeInput';

export const CreateTestFormAuthInput: React.FC<{form: FormInstance<ICreateTestValues>}> = ({form}) => {
  return (
    <>
      <Divider />
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
    </>
  );
};
