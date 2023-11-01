import {Form} from 'antd';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from '../RequestDetails.styled';
import * as R from './RequestDetailsAuthInput.styled';

interface IProps {
  baseName: string[];
}

export const ApiKeyFieldsBase = ({baseName}: IProps) => (
  <S.Row>
    <R.FlexContainer>
      <Form.Item
        data-cy="apiKey-key"
        style={{flexBasis: '50%'}}
        name={[...baseName, 'apiKey', 'key']}
        label="Key"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder="Enter key" />
      </Form.Item>
      <Form.Item
        data-cy="apiKey-value"
        style={{flexBasis: '50%'}}
        name={[...baseName, 'apiKey', 'value']}
        label="Value"
        rules={[{required: true}]}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder="Enter value" />
      </Form.Item>
    </R.FlexContainer>
  </S.Row>
);
