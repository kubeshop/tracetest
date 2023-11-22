import {Form, FormInstance, Input} from 'antd';

import {Headers} from 'components/Fields';
import {DEFAULT_VALUES} from 'components/VariableSetModal/VariableSetModal';
import VariableSet from 'models/VariableSet.model';

interface IProps {
  form: FormInstance<VariableSet>;
  initialValues?: VariableSet;
  onSubmit(values: VariableSet): void;
  onValidate(changedValues: any, values: VariableSet): void;
}

const VariableSetForm = ({form, initialValues, onSubmit, onValidate}: IProps) => {
  return (
    <Form<VariableSet>
      initialValues={{...initialValues}}
      form={form}
      layout="vertical"
      name="variableSet"
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
      >
        <Input />
      </Form.Item>

      <Headers
        initialValue={form.getFieldValue('values') || DEFAULT_VALUES}
        label="Values"
        name={['values']}
        unit="Key"
      />
    </Form>
  );
};

export default VariableSetForm;
