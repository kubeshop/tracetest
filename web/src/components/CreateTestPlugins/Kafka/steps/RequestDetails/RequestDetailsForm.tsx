import {Form, Tabs} from 'antd';
import KeyValueListInput from 'components/Fields/KeyValueList';
import {PlainAuth, SSL, SkipTraceCollection} from 'components/Fields';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from './RequestDetails.styled';

const RequestDetailsForm = () => {
  return (
    <Tabs defaultActiveKey="auth">
      <Tabs.TabPane forceRender tab="Auth" key="auth">
        <PlainAuth />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Message" key="message">
        <Form.Item
          label="Topic"
          data-cy="topic"
          name="topic"
          rules={[{required: true, message: 'Please enter a topic'}]}
        >
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
        <KeyValueListInput
          name="headers"
          label="Message Headers"
          addButtonLabel="Add Header"
          keyPlaceholder="Header Key"
          valuePlaceholder="Header Value"
          initialValue={[{key: '', value: ''}]}
        />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Settings" key="settings">
        <S.SettingsContainer>
          <SSL />
          <SkipTraceCollection />
        </S.SettingsContainer>
      </Tabs.TabPane>
    </Tabs>
  );
};

export default RequestDetailsForm;
