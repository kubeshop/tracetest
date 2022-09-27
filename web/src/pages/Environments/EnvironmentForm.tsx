import {Button, Form, FormInstance, Input} from 'antd';
import {useEffect} from 'react';
import styled from 'styled-components';
import RequestDetailsHeadersInput from '../../components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import {
  IEnvironment,
  useCreateEnvironmentMutation,
  useLazyGetEnvironmentSecretListQuery,
} from '../../redux/apis/TraceTest.api';
import {EnvironmentState} from './EnvironmentState';

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
`;

interface IProps {
  form: FormInstance<IEnvironment>;
  state: EnvironmentState;
  onCancel: () => void;
}

export const EnvironmentForm: React.FC<IProps> = ({state, form, onCancel}) => {
  const [createEnvironment] = useCreateEnvironmentMutation();
  const [loadResultList, {data}] = useLazyGetEnvironmentSecretListQuery();
  useEffect(() => {
    if (state.environment) loadResultList({environmentId: state.environment?.id || ''});
  }, [loadResultList, state.environment]);
  useEffect(() => form.setFieldsValue({variables: data || []}), [form, data]);
  useEffect(() => {
    form.setFieldsValue({name: state.environment?.name, description: state.environment?.description});
  }, [form, state.environment?.name, state.environment?.description]);
  return (
    <Form<IEnvironment>
      name="basic"
      layout="vertical"
      form={form}
      onFinish={s => {
        createEnvironment({...s});
        onCancel();
      }}
    >
      <Form.Item label="Name" name="name" rules={[{required: true, message: 'Please input your name!'}]}>
        <Input />
      </Form.Item>
      <Form.Item
        label="Description"
        name="description"
        rules={[{required: true, message: 'Please input your description!'}]}
      >
        <Input />
      </Form.Item>
      <RequestDetailsHeadersInput name="variables" unit="Key" initialValue={undefined} addLabel="Add Entry" />
      <Footer>
        <Button style={{marginRight: 16}} data-cy="create-test-cancel" type="primary" ghost onClick={onCancel}>
          Cancel
        </Button>
        <Form.Item shouldUpdate className="submit">
          {() => (
            <Button
              type="primary"
              htmlType="submit"
              disabled={
                !form.isFieldsTouched(true) || form.getFieldsError().filter(({errors}: any) => errors.length).length > 0
              }
            >
              Log in
            </Button>
          )}
        </Form.Item>
      </Footer>
    </Form>
  );
};
