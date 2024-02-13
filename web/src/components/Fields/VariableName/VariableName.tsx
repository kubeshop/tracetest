import {Form} from 'antd';
import SingleLine from '../../Inputs/SingleLine';

const VariableName = () => (
  <Form.Item
    name="id"
    rules={[{required: true, message: 'Please enter a valid variable name'}]}
    style={{marginBottom: 0}}
  >
    <SingleLine placeholder="Enter variable name" />
  </Form.Item>
);

export default VariableName;
