import {Form} from 'antd';

import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';

const VariableNameForm = () => {
  return (
    <S.InputContainer>
      <Form.Item label="Variable Name" name="id" rules={[{required: true, message: 'Please enter a Variable Name'}]}>
        <Editor type={SupportedEditors.Interpolation} placeholder="Variable Name" />
      </Form.Item>
    </S.InputContainer>
  );
};

export default VariableNameForm;
