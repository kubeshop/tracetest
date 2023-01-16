import {Form} from 'antd';

import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';

const ValueForm = () => {
  return (
    <Form.Item
      label="Trace ID"
      name="id"
      rules={[{required: true, message: 'Please enter a trace id or an expression'}]}
    >
      <Editor type={SupportedEditors.Interpolation} placeholder="Trace id or expression" />
    </Form.Item>
  );
};

export default ValueForm;
