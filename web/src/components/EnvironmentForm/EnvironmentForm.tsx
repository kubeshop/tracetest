import {Form, FormInstance, Input} from 'antd';

import RequestDetailsHeadersInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import {DEFAULT_VALUES} from 'pages/Environments/EnvironmentModal';
import {TEnvironment} from 'types/Environment.types';

interface IProps {
  form: FormInstance<TEnvironment>;
  initialValues?: TEnvironment;
  onSubmit(values: TEnvironment): void;
  onValidate(changedValues: any, values: TEnvironment): void;
}

const EnvironmentForm = ({form, initialValues, onSubmit, onValidate}: IProps) => {
  return (
    <Form<TEnvironment>
      initialValues={{...initialValues}}
      form={form}
      layout="vertical"
      name="environment"
      onFinish={onSubmit}
      onValuesChange={onValidate}
    >
      <Form.Item label="Name" name="name" rules={[{required: true, message: 'Please input a name'}]}>
        <Input />
      </Form.Item>

      <Form.Item
        label="Description"
        name="description"
        rules={[{required: true, message: 'Please input a description'}]}
      >
        <Input />
      </Form.Item>

      <RequestDetailsHeadersInput
        initialValue={form.getFieldValue('values') || DEFAULT_VALUES}
        label="Values"
        name="values"
        unit="Key"
      />
    </Form>
  );
};

export default EnvironmentForm;
