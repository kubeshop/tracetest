import {Form} from 'antd';
import React from 'react';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from '../RequestDetails.styled';
import * as R from './RequestDetailsAuthInput.styled';

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
        <Editor type={SupportedEditors.Interpolation} />
      </Form.Item>
      <Form.Item
        style={{flexBasis: '50%'}}
        name={['auth', 'basic', 'password']}
        label="Password"
        data-cy="basic-password"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} />
      </Form.Item>
    </R.FlexContainer>
  </S.Row>
);
