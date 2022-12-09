import {Form, Select} from 'antd';
import {ApiKeyFieldsBase} from './apiKeyFieldsBase';

interface IProps {
  baseName: string[];
}

export const ApiKeyFields = ({baseName}: IProps) => (
  <>
    <ApiKeyFieldsBase baseName={baseName} />
    <Form.Item style={{minWidth: '100%'}} initialValue="query" label="Add To" name={[...baseName, 'apiKey', 'in']}>
      <Select
        className="select-auth-method"
        data-cy="auth-apiKey-select"
        dropdownClassName="select-dropdown-auth-method"
      >
        {['query', 'header'].map(m => (
          <Select.Option data-cy={`auth-apiKey-select-option-${m}`} key={m} value={m}>
            {m === 'query' ? 'Query Params' : 'Header'}
          </Select.Option>
        ))}
      </Select>
    </Form.Item>
  </>
);
