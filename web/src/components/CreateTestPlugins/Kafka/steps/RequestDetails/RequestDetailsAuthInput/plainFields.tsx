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
        style={{flexBasis: '49%', marginTop: '26px'}}
        name={[...baseName, 'basic', 'username']}
        data-cy="basic-username"
        label="Username"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} />
      </Form.Item>
      <Form.Item
        style={{flexBasis: '2%', marginTop: '26px'}}>
        <div />
      </Form.Item>
      <Form.Item
        style={{flexBasis: '49%', marginTop: '26px'}}
        name={[...baseName, 'basic', 'password']}
        label="Password"
        data-cy="basic-password"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} />
      </Form.Item>
    </R.FlexContainer>
  </S.Row>
);
