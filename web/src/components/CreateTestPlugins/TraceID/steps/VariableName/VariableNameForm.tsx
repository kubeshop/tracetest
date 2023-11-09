import {Form} from 'antd';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';

const VariableNameForm = () => (
  <Form.Item name="id" rules={[{required: true, message: 'Please enter a Variable Name'}]}>
    <Editor type={SupportedEditors.Interpolation} placeholder="Variable Name" />
  </Form.Item>
);

export default VariableNameForm;
