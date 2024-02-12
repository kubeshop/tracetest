import {Form} from 'antd';
import * as S from './Auth.styled';
import SingleLine from '../../Inputs/SingleLine';

interface IProps {
  baseName: string[];
}

const AuthBasic = ({baseName}: IProps) => (
  <S.Row>
    <S.FlexContainer>
      <Form.Item
        style={{flexBasis: '50%', overflow: 'hidden'}}
        name={[...baseName, 'basic', 'username']}
        data-cy="basic-username"
        label="Username"
        rules={[{required: true}]}
      >
        <SingleLine />
      </Form.Item>
      <Form.Item
        style={{flexBasis: '50%', overflow: 'hidden'}}
        name={[...baseName, 'basic', 'password']}
        label="Password"
        data-cy="basic-password"
        rules={[{required: true}]}
      >
        <SingleLine />
      </Form.Item>
    </S.FlexContainer>
  </S.Row>
);

export default AuthBasic;
