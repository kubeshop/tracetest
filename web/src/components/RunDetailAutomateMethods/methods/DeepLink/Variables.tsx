import {PlusOutlined} from '@ant-design/icons';
import {Form, Input, Typography} from 'antd';
import * as S from './DeepLink.styled';

const Variables = () => (
  <S.VariablesContainer>
    <Form.List name="variables">
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            // eslint-disable-next-line react/no-array-index-key
            <S.EntryContainer key={`${field.name}-${index}`}>
              <S.ValuesContainer>
                <label htmlFor={`${field.name}-${index}-key`}>
                  <Typography.Text>Variable Key</Typography.Text>
                </label>
                <label htmlFor={`${field.name}-${index}-value`}>
                  <Typography.Text>Variable Value</Typography.Text>
                </label>
                <div />
              </S.ValuesContainer>
              <S.ValuesContainer>
                <div>
                  <Form.Item name={[field.name, 'key']} noStyle>
                    <Input placeholder="variable key" id={`${field.name}-${index}-key`} />
                  </Form.Item>
                </div>
                <div>
                  <Form.Item name={[field.name, 'value']} noStyle>
                    <Input placeholder="variable value" id={`${field.name}-${index}-value`} />
                  </Form.Item>
                </div>
                <S.DeleteVariableButton icon={<S.DeleteIcon />} onClick={() => remove(field.name)} type="link" />
              </S.ValuesContainer>
            </S.EntryContainer>
          ))}

          <S.AddVariableButton data-cy="add-header" icon={<PlusOutlined />} onClick={() => add()} type="link">
            Add Variable
          </S.AddVariableButton>
        </>
      )}
    </Form.List>
  </S.VariablesContainer>
);

export default Variables;
