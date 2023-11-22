import {Form, Tabs} from 'antd';
import {KeyValueList, PlainAuth, SSL, SkipTraceCollection} from 'components/Fields';
import {Editor} from 'components/Inputs';
import useQueryTabs from 'hooks/useQueryTabs';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from './Kafka.styled';

const Kafka = () => {
  const [activeKey, setActiveKey] = useQueryTabs('auth', 'triggerTab');

  return (
    <Tabs defaultActiveKey={activeKey} onChange={setActiveKey} activeKey={activeKey}>
      <Tabs.TabPane forceRender tab="Auth" key="auth">
        <PlainAuth />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Message" key="message">
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

      <Tabs.TabPane forceRender tab="Topic" key="topic">
        <Form.Item data-cy="topic" name="topic" rules={[{required: true, message: 'Please enter a topic'}]}>
          <Editor type={SupportedEditors.Interpolation} placeholder="my-topic" />
        </Form.Item>
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Headers" key="headers">
        <KeyValueList
          name="headers"
          label=""
          addButtonLabel="Add Header"
          keyPlaceholder="Header"
          valuePlaceholder="Value"
          initialValue={[{key: '', value: ''}]}
        />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab="Settings" key="settings">
        <S.SettingsContainer>
          <SSL formID="kafka" />
          <SkipTraceCollection />
        </S.SettingsContainer>
      </Tabs.TabPane>
    </Tabs>
  );
};

export default Kafka;
