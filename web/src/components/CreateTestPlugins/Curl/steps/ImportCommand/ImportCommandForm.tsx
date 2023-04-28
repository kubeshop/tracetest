import {Form} from 'antd';
import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import CurlService from 'services/Triggers/Curl.service';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';

export const FORM_ID = 'create-test';

const ImportCommandForm = () => {
  return (
    <S.InputContainer>
      <Form.Item
        label="Paste the CURL command"
        name="command"
        rules={[
          {required: true, message: 'Please enter a command'},
          {
            validator: (_, command) => {
              if (!CurlService.getIsValidCommand(command)) {
                return Promise.reject(new Error('Invalid CURL command'));
              }

              return Promise.resolve();
            },
            message: 'Invalid CURL command',
          },
        ]}
        style={{marginBottom: 0}}
      >
        <Editor type={SupportedEditors.CurlCommand} />
      </Form.Item>
    </S.InputContainer>
  );
};

export default ImportCommandForm;
