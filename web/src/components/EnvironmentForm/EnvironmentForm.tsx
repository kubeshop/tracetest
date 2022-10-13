import {Button, Form, FormInstance, Input} from 'antd';
import RequestDetailsHeadersInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import {useEffect, useMemo} from 'react';
import {TEnvironment} from '../../types/Environment.types';
import {useCreateEnvironmentMutation, useLazyGetEnvironmentSecretListQuery} from '../../redux/apis/TraceTest.api';
import {Footer} from './EnvironmentForm.styled';

interface IProps {
  form: FormInstance<TEnvironment>;
  environment?: TEnvironment;
  onCancel: () => void;
}

const EnvironmentForm: React.FC<IProps> = ({environment, form, onCancel}) => {
  const [createEnvironment] = useCreateEnvironmentMutation();
  const [loadResultList, {data}] = useLazyGetEnvironmentSecretListQuery();
  const isDisabled = useMemo(
    () => !form.isFieldsTouched(true) || form.getFieldsError().filter(({errors}: any) => errors.length).length > 0,
    [form]
  );
  useEffect(() => {
    if (environment) loadResultList({environmentId: environment?.id || ''});
  }, [loadResultList, environment]);
  useEffect(() => form.setFieldsValue({variables: data || []}), [form, data]);
  useEffect(() => {
    form.setFieldsValue({name: environment?.name, description: environment?.description});
  }, [form, environment?.name, environment?.description]);
  return (
    <Form<TEnvironment>
      name="basic"
      layout="vertical"
      form={form}
      onFinish={s => {
        createEnvironment(s);
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
      <RequestDetailsHeadersInput name="variables" unit="Key" initialValue={undefined} label="Entry" />
      <Footer>
        <Button style={{marginRight: 16}} data-cy="create-test-cancel" type="primary" ghost onClick={onCancel}>
          Cancel
        </Button>
        <Form.Item shouldUpdate className="submit">
          <Button type="primary" htmlType="submit" disabled={isDisabled}>
            Create
          </Button>
        </Form.Item>
      </Footer>
    </Form>
  );
};

export default EnvironmentForm;
