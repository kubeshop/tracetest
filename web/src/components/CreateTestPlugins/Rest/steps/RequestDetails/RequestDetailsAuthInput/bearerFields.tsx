import {Form} from 'antd';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';

export const bearerFields: React.ReactElement = (
  <Form.Item data-cy="bearer-token" name={['auth', 'bearer', 'token']} label="Token" rules={[{required: true}]}>
    <Editor type={SupportedEditors.Interpolation} />
  </Form.Item>
);
