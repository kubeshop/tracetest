import {Form} from 'antd';
import Fields from './Fields';
import PlainAuthType from './PlainAuthType';

interface IProps {
  name?: string[];
}

const PlainAuth = ({name = ['authentication']}: IProps) => (
  <>
    <PlainAuthType baseName={name} />
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
  </>
);

export default PlainAuth;
