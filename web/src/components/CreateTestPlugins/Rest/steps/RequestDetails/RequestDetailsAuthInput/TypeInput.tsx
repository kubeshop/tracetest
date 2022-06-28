import {Form, FormInstance, Select} from 'antd';
import {IRequestDetailsValues} from '../RequestDetails';
import * as S from '../RequestDetails.styled';

const methodNamingMap: Record<string, string> = {
  apiKey: 'API Key',
  bearer: 'Bearer Token',
  basic: 'Basic Auth',
  none: 'No Auth',
};

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
        {[null, 'apiKey', 'basic', 'bearer'].map(method => (
          <Select.Option data-cy={`auth-type-select-option-${method}`} key={method} value={method}>
            {method ? methodNamingMap[method] : methodNamingMap.none}
          </Select.Option>
        ))}
      </Select>
    </Form.Item>
  </S.Row>
);

export default TypeInput;
