import {Form} from 'antd';
import * as S from './PlainAuth.styled';
import SingleLine from '../../Inputs/SingleLine';

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
        style={{flexBasis: '50%', overflow: 'hidden'}}
      >
        <SingleLine placeholder="Kafka Plain Username" />
      </Form.Item>

      <Form.Item
        name={[...baseName, 'plain', 'password']}
        label="Password"
        data-cy="plain-password"
        rules={[{required: true}]}
        style={{flexBasis: '50%', overflow: 'hidden'}}
      >
        <SingleLine placeholder="Kafka Plain Password" />
      </Form.Item>
    </S.FlexContainer>
  </S.Row>
);

export default Fields;
