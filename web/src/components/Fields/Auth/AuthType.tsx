import {Form, Select} from 'antd';
import * as S from '../../TestPlugins/Forms/Rest/Rest.styled';

const authMethodList = [
  {
    name: 'No Auth',
    value: null,
  },
  {
    name: 'API Key',
    value: 'apiKey',
  },
  {
    name: 'Bearer Token',
    value: 'bearer',
  },
  {
    name: 'Basic Auth',
    value: 'basic',
  },
] as const;

interface IProps {
  baseName: string[];
}

const AuthType = ({baseName}: IProps) => (
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

export default AuthType;
