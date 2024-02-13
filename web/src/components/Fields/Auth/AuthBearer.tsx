import {Form} from 'antd';
import SingleLine from '../../Inputs/SingleLine';

interface IProps {
  baseName: string[];
}

const AuthBearer = ({baseName}: IProps) => (
  <Form.Item data-cy="bearer-token" name={[...baseName, 'bearer', 'token']} label="Token" rules={[{required: true}]}>
    <SingleLine />
  </Form.Item>
);

export default AuthBearer;
