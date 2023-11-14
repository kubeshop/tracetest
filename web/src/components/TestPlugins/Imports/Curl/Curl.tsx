import {Form} from 'antd';

import CurlService from 'services/Importers/Curl.service';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';

export const FORM_ID = 'create-test';

const Curl = () => {
  return (
    <Form.Item
      label="Paste cUrl command, raw text, or url to definition file"
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
    >
      <Editor placeholder="curl -X GET http://localhost:8080/" type={SupportedEditors.CurlCommand} />
    </Form.Item>
  );
};

export default Curl;
