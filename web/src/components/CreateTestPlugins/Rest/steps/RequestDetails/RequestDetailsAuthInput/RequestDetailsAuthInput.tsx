import {Form} from 'antd';
import {ApiKeyFields} from './apiKeyFields';
import {ApiKeyFieldsBase} from './apiKeyFieldsBase';
import {BasicFields} from './basicFields';
import {BearerFields} from './bearerFields';
import TypeInput from './TypeInput';

interface IProps {
  hasBaseApikeyFields?: boolean;
  name?: string[];
}

const RequestDetailsAuthInput = ({hasBaseApikeyFields = false, name = ['auth']}: IProps) => (
  <div>
    <TypeInput baseName={name} />
    <Form.Item noStyle shouldUpdate style={{marginBottom: 0, width: '100%'}}>
      {({getFieldValue}) => {
        const method = getFieldValue(name)?.type;
        switch (method) {
          case 'bearer':
            return <BearerFields baseName={name} />;
          case 'basic':
            return <BasicFields baseName={name} />;
          case 'apiKey':
            return hasBaseApikeyFields ? <ApiKeyFieldsBase baseName={name} /> : <ApiKeyFields baseName={name} />;
          default:
            return null;
        }
      }}
    </Form.Item>
  </div>
);

export default RequestDetailsAuthInput;
