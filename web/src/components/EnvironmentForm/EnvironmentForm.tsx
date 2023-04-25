import {Form, FormInstance, Input} from 'antd';

import RequestDetailsHeadersInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import {DEFAULT_VALUES} from 'components/EnvironmentModal/EnvironmentModal';
import Environment from 'models/Environment.model';

interface IProps {
  form: FormInstance<Environment>;
  initialValues?: Environment;
  onSubmit(values: Environment): void;
  onValidate(changedValues: any, values: Environment): void;
}

const EnvironmentForm = ({form, initialValues, onSubmit, onValidate}: IProps) => {
  return (
    <Form<Environment>
      initialValues={{...initialValues}}
      form={form}
      layout="vertical"
      name="environment"
      onFinish={onSubmit}
      onValuesChange={onValidate}
    >
      <Form.Item name="id" hidden>
        <Input type="hidden" />
      </Form.Item>
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
        name={['values']}
        unit="Key"
      />
    </Form>
  );
};

export default EnvironmentForm;
