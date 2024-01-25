import {Form} from 'antd';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from './Auth.styled';

interface IProps {
  baseName: string[];
}

const AuthApiKeyBase = ({baseName}: IProps) => (
  <S.Row>
    <S.FlexContainer>
      <Form.Item
        data-cy="apiKey-key"
        style={{flexBasis: '50%', overflow: 'hidden'}}
        name={[...baseName, 'apiKey', 'key']}
        label="Key"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder="Enter key" />
      </Form.Item>
      <Form.Item
        data-cy="apiKey-value"
        style={{flexBasis: '50%', overflow: 'hidden'}}
        name={[...baseName, 'apiKey', 'value']}
        label="Value"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder="Enter value" />
      </Form.Item>
    </S.FlexContainer>
  </S.Row>
);

export default AuthApiKeyBase;
