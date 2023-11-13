import {Form} from 'antd';
import React from 'react';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from './Auth.styled';

interface IProps {
  baseName: string[];
}

const AuthBasic = ({baseName}: IProps) => (
  <S.Row>
    <S.FlexContainer>
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
    </S.FlexContainer>
  </S.Row>
);

export default AuthBasic;
