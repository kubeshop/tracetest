import {Form, Input, Select} from 'antd';
import * as S from '../CreateTestModal.styled';
import * as R from './CreateTestFormAuthInput.styled';

export const apiKeyFields: React.ReactElement = (
  <>
    <S.Row>
      <R.FlexContainer>
        <Form.Item
          data-cy="apiKey-key"
          style={{flexBasis: '50%'}}
          name={['auth', 'apiKey', 'key']}
          label="Key"
          rules={[{required: true}]}
        >
          <Input placeholder="Enter key" />
        </Form.Item>
        <Form.Item
          data-cy="apiKey-value"
          style={{flexBasis: '50%'}}
          name={['auth', 'apiKey', 'value']}
          label="Value"
          rules={[{required: true}]}
        >
          <Input placeholder="Enter value" />
        </Form.Item>
      </R.FlexContainer>
    </S.Row>
    <Form.Item style={{minWidth: '100%'}} initialValue="query" label="Add To" name={['auth', 'apiKey', 'in']}>
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
