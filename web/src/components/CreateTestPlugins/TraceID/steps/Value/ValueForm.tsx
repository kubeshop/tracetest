import {Form} from 'antd';

import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';

const ValueForm = () => {
  return (
    <S.InputContainer>
      <Form.Item
        label="Trace ID"
        name="id"
        rules={[{required: true, message: 'Please enter a Trace ID or an Expression'}]}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder="Trace ID or Expression" />
      </Form.Item>
    </S.InputContainer>
  );
};

export default ValueForm;
