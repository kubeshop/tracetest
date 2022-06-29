import {Form, FormInstance, Select} from 'antd';
import {IRequestDetailsValues} from '../RequestDetails';
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
  form: FormInstance<IRequestDetailsValues>;
}

const TypeInput = ({form}: IProps) => (
  <S.Row>
    <Form.Item
      style={{minWidth: '100%'}}
      initialValue={null}
      label="Authorization Type"
      name={['auth', 'type']}
      valuePropName="type"
    >
      <Select
        className="select-auth-method"
        data-cy="auth-type-select"
        dropdownClassName="select-dropdown-auth-type"
        placeholder="No Auth"
        allowClear
        onClear={() => form.resetFields(['auth'])}
        onChange={e => {
          form.resetFields(['auth']);
          form.setFieldsValue({auth: {type: e as any}});
        }}
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
