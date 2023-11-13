import {Form} from 'antd';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';

interface IProps {
  baseName: string[];
}

const AuthBearer = ({baseName}: IProps) => (
  <Form.Item data-cy="bearer-token" name={[...baseName, 'bearer', 'token']} label="Token" rules={[{required: true}]}>
    <Editor type={SupportedEditors.Interpolation} />
  </Form.Item>
);

export default AuthBearer;
