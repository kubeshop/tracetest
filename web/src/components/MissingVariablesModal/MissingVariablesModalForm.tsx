import {Form} from 'antd';
import TestListVariablesInput from './TestVariablesInput/TestListVariablesInput';

const MissingVariablesModalForm = () => {
  return (
    <Form.Item name="variables">
      <TestListVariablesInput />
    </Form.Item>
  );
};

export default MissingVariablesModalForm;
