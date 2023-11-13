import {Form} from 'antd';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from './PlainAuth.styled';

interface IProps {
  baseName: string[];
}

const Fields = ({baseName}: IProps) => (
  <S.Row>
    <S.FlexContainer>
      <S.PlainFieldUsername>
        <Form.Item
          name={[...baseName, 'plain', 'username']}
          data-cy="plain-username"
          label="Username"
          rules={[{required: true}]}
        >
          <Editor type={SupportedEditors.Interpolation} placeholder="Kafka Plain Username" />
        </Form.Item>
      </S.PlainFieldUsername>
      <S.PlainFieldPassword>
        <Form.Item
          name={[...baseName, 'plain', 'password']}
          label="Password"
          data-cy="plain-password"
          rules={[{required: true}]}
        >
          <Editor type={SupportedEditors.Interpolation} placeholder="Kafka Plain Password" />
        </Form.Item>
      </S.PlainFieldPassword>
    </S.FlexContainer>
  </S.Row>
);

export default Fields;
