import {Form} from 'antd';
import React from 'react';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from '../RequestDetails.styled';
import * as R from './RequestDetailsAuthInput.styled';

interface IProps {
  baseName: string[];
}

export const BasicFields = ({baseName}: IProps) => (
  <S.Row>
    <R.FlexContainer>
      <Form.Item
        style={{flexBasis: '50%'}}
        name={[...baseName, 'basic', 'username']}
        data-cy="basic-username"
        label="Username"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} />
      </Form.Item>
      <Form.Item
        style={{flexBasis: '50%'}}
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
