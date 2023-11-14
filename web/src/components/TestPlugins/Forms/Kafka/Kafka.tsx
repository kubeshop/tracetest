import {Form, Tabs} from 'antd';
import {KeyValueList, PlainAuth, SSL} from 'components/Fields';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import {FORM_ID} from 'pages/CreateTest/Content';

const RequestDetailsForm = () => (
  <Tabs defaultActiveKey="auth">
    <Tabs.TabPane forceRender tab="Auth" key="auth">
      <PlainAuth />
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Message" key="message">
      <Form.Item label="Topic" data-cy="topic" name="topic" rules={[{required: true, message: 'Please enter a topic'}]}>
        <Editor type={SupportedEditors.Interpolation} placeholder="my-topic" />
      </Form.Item>

      <Form.Item label="Key" data-cy="message-key" name="messageKey">
        <Editor type={SupportedEditors.Interpolation} placeholder="my-message-name" />
      </Form.Item>

      <Form.Item
        label="Value"
        data-cy="message-value"
        name="messageValue"
        rules={[{required: true, message: 'Please enter a message value'}]}
      >
        <Editor type={SupportedEditors.Interpolation} placeholder="my-message-value" />
      </Form.Item>
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Headers" key="headers">
      <KeyValueList
        name="headers"
        label="Message Headers"
        addButtonLabel="Add Header"
        keyPlaceholder="Header Key"
        valuePlaceholder="Header Value"
        initialValue={[{key: '', value: ''}]}
      />
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Settings" key="settings">
      <SSL formID={FORM_ID} />
    </Tabs.TabPane>
  </Tabs>
);

export default RequestDetailsForm;
