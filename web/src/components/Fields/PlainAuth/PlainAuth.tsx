import {Form} from 'antd';
import Fields from './Fields';
import TypeInput from './Input';

interface IProps {
  name?: string[];
}

const PlainAuth = ({name = ['authentication']}: IProps) => (
  <div>
    <TypeInput baseName={name} />
    <Form.Item noStyle shouldUpdate style={{marginBottom: 0, width: '100%'}}>
      {({getFieldValue}) => {
        const method = getFieldValue(name)?.type;
        switch (method) {
          case 'plain':
            return <Fields baseName={name} />;
          default:
            return null;
        }
      }}
    </Form.Item>
  </div>
);

export default PlainAuth;
