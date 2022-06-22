import {Form, Input} from 'antd';
import React from 'react';
import * as S from '../RequestDetails.styled';
import * as R from './CreateTestFormAuthInput.styled';

export const basicFields: React.ReactElement = (
  <S.Row>
    <R.FlexContainer>
      <Form.Item
        style={{flexBasis: '50%'}}
        name={['auth', 'basic', 'username']}
        data-cy="basic-username"
        label="Username"
        rules={[{required: true}]}
      >
        <Input />
      </Form.Item>
      <Form.Item
        style={{flexBasis: '50%'}}
        name={['auth', 'basic', 'password']}
        label="Password"
        data-cy="basic-password"
        rules={[{required: true}]}
      >
        <Input />
      </Form.Item>
    </R.FlexContainer>
  </S.Row>
);
