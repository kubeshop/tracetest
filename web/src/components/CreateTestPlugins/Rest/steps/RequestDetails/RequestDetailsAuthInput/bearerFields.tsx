import {Form} from 'antd';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';

interface IProps {
  baseName: string[];
}

export const BearerFields = ({baseName}: IProps) => (
  <Form.Item data-cy="bearer-token" name={[...baseName, 'bearer', 'token']} label="Token" rules={[{required: true}]}>
    <Editor type={SupportedEditors.Interpolation} />
  </Form.Item>
);
