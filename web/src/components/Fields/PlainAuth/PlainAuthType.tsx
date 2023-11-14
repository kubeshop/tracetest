import {Form, Select} from 'antd';
import * as S from '../../TestPlugins/Forms/Kafka/Kafka.styled';

const authMethodList = [
  {
    name: 'No Auth',
    value: null,
  },
  {
    name: 'Plain',
    value: 'plain',
  },
] as const;

interface IProps {
  baseName: string[];
}

const PlainAuthType = ({baseName}: IProps) => (
  <S.Row>
    <Form.Item shouldUpdate style={{minWidth: '100%'}} name={[...baseName, 'type']}>
      <Select data-cy="auth-type-select" defaultValue={null}>
        {authMethodList.map(({name, value}) => (
          <Select.Option data-cy={`auth-type-select-option-${value}`} key={value} value={value}>
            {name}
          </Select.Option>
        ))}
      </Select>
    </Form.Item>
  </S.Row>
);

export default PlainAuthType;
