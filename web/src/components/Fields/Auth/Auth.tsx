import {Form} from 'antd';
import AuthApiKey from './AuthApiKey';
import AuthApiKeyBase from './AuthApiKeyBase';
import AuthBasic from './AuthBasic';
import AuthBearer from './AuthBearer';
import TypeInput from './AuthType';

interface IProps {
  hasBaseApikeyFields?: boolean;
  name?: string[];
}

const Auth = ({hasBaseApikeyFields = false, name = ['auth']}: IProps) => (
  <>
    <TypeInput baseName={name} />
    <Form.Item noStyle shouldUpdate style={{marginBottom: 0, width: '100%'}}>
      {({getFieldValue}) => {
        const method = getFieldValue(name)?.type;
        switch (method) {
          case 'bearer':
            return <AuthBearer baseName={name} />;
          case 'basic':
            return <AuthBasic baseName={name} />;
          case 'apiKey':
            return hasBaseApikeyFields ? <AuthApiKeyBase baseName={name} /> : <AuthApiKey baseName={name} />;
          default:
            return null;
        }
      }}
    </Form.Item>
  </>
);

export default Auth;
