import {Form} from 'antd';
import {PlainUserPasswordFields} from './PlainUserPasswordFields';
import TypeInput from './TypeInput';

interface IProps {
  name?: string[];
}

const RequestDetailsAuthInput = ({name = ['authentication']}: IProps) => (
  <div>
    <TypeInput baseName={name} />
    <Form.Item noStyle shouldUpdate style={{marginBottom: 0, width: '100%'}}>
      {({getFieldValue}) => {
        const method = getFieldValue(name)?.type;
        switch (method) {
          case 'plain':
            return <PlainUserPasswordFields baseName={name} />;
          default:
            return null;
        }
      }}
    </Form.Item>
  </div>
);

export default RequestDetailsAuthInput;
