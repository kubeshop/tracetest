import {Form} from 'antd';
import React from 'react';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from '../RequestDetails.styled';
import * as R from './RequestDetailsAuthInput.styled';

interface IProps {
  baseName: string[];
}

export const PlainFields = ({baseName}: IProps) => (
  <S.Row>
    <R.FlexContainer>
      <Form.Item
        style={{flexBasis: '49%', marginTop: '26px', marginRight: '2px'}}
        name={[...baseName, 'plain', 'username']}
        data-cy="plain-username"
        label="Username"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder='Kafka Plain Username' />
      </Form.Item>
      <Form.Item
        style={{flexBasis: '49%', marginTop: '26px'}}
        name={[...baseName, 'plain', 'password']}
        label="Password"
        data-cy="plain-password"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder='Kafka Plain Password' />
      </Form.Item>
    </R.FlexContainer>
  </S.Row>
);
