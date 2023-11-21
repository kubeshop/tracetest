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
      <Form.Item
        name={[...baseName, 'plain', 'username']}
        data-cy="plain-username"
        label="Username"
        rules={[{required: true}]}
        style={{flexBasis: '50%'}}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder="Kafka Plain Username" />
      </Form.Item>

      <Form.Item
        name={[...baseName, 'plain', 'password']}
        label="Password"
        data-cy="plain-password"
        rules={[{required: true}]}
        style={{flexBasis: '50%'}}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder="Kafka Plain Password" />
      </Form.Item>
    </S.FlexContainer>
  </S.Row>
);

export default Fields;
