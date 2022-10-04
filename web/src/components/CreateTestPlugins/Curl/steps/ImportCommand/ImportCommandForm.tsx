import CodeMirror from '@uiw/react-codemirror';
import {StreamLanguage} from '@codemirror/language';
import {shell} from '@codemirror/legacy-modes/mode/shell';
import {Form} from 'antd';
import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import CurlService from 'services/Triggers/Curl.service';

export const FORM_ID = 'create-test';

const ImportCommandForm = () => {
  return (
    <S.InputContainer>
      <Form.Item
        name="command"
        rules={[
          {required: true, message: 'Please enter a command'},
          {
            validator: (_, command) => {
              if (!CurlService.getIsValidCommand(command)) throw new Error('Invalid command');

              return Promise.resolve(true);
            },
            message: 'Invalid CURL command',
          },
        ]}
        style={{marginBottom: 0}}
      >
        <CodeMirror
          data-cy="body"
          basicSetup={{lineNumbers: true, indentOnInput: true}}
          extensions={[StreamLanguage.define(shell)]}
          spellCheck={false}
          placeholder="Enter a curl command"
        />
      </Form.Item>
    </S.InputContainer>
  );
};

export default ImportCommandForm;
