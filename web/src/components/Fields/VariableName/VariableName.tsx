import {Form} from 'antd';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';

const VariableName = () => (
  <Form.Item
    name="id"
    rules={[{required: true, message: 'Please enter a valid variable name'}]}
    style={{marginBottom: 0}}
  >
    <Editor type={SupportedEditors.Interpolation} placeholder="Enter variable name" />
  </Form.Item>
);

export default VariableName;
