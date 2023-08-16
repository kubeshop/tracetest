import {Form, Select} from 'antd';
import * as S from '../RequestDetails.styled';

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

const TypeInput = ({baseName}: IProps) => (
  <S.Row>
    <Form.Item shouldUpdate style={{minWidth: '100%'}} label="Authentication Type" name={[...baseName, 'type']}>
      <Select
        className="select-auth-method"
        data-cy="auth-type-select"
        dropdownClassName="select-dropdown-auth-type"
        placeholder="No Auth"
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
