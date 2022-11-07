import {Form, FormInstance, Select} from 'antd';
import {TRequestAuth} from 'types/Test.types';
import * as S from '../RequestDetails.styled';

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
  form: FormInstance<Partial<{auth: TRequestAuth}>>;
}

const TypeInput = ({form}: IProps) => (
  <S.Row>
    <Form.Item shouldUpdate style={{minWidth: '100%'}} label="Authorization Type" name={['auth', 'type']}>
      <Select
        className="select-auth-method"
        data-cy="auth-type-select"
        dropdownClassName="select-dropdown-auth-type"
        placeholder="No Auth"
        allowClear
        onClear={() => form.resetFields(['auth'])}
      >
        {authMethodList.map(({name, value}) => (
          <Select.Option data-cy={`auth-type-select-option-${value}`} key={value} value={value}>
            {name}
          </Select.Option>
        ))}
      </Select>
    </Form.Item>
  </S.Row>
);

export default TypeInput;
