import {Form, Tabs} from 'antd';
import {KeyValueList, PlainAuth, SSL, SkipTraceCollection} from 'components/Fields';
import {Editor} from 'components/Inputs';
import TriggerTab from 'components/TriggerTab';
import useQueryTabs from 'hooks/useQueryTabs';
import {SupportedEditors} from 'constants/Editor.constants';
import {TDraftTest} from 'types/Test.types';
import * as S from './Kafka.styled';
import SingleLine from '../../../Inputs/SingleLine';

const Kafka = () => {
  const [activeKey, setActiveKey] = useQueryTabs('auth', 'triggerTab');
  const form = Form.useFormInstance<TDraftTest>();
  const authType = Form.useWatch(['authentication', 'type'], form);
  const messageValue = Form.useWatch('messageValue', form);
  const topic = Form.useWatch('topic', form);
  const headers = Form.useWatch('headers', form);
  const sslVerification = Form.useWatch('sslVerification', form);
  const skipTraceCollection = Form.useWatch('skipTraceCollection', form);

  return (
    <Tabs defaultActiveKey={activeKey} onChange={setActiveKey} activeKey={activeKey}>
      <Tabs.TabPane forceRender tab={<TriggerTab hasContent={!!authType} label="Auth" />} key="auth">
        <PlainAuth />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab={<TriggerTab hasContent={!!messageValue} label="Message" />} key="message">
        <Form.Item label="Key" data-cy="message-key" name="messageKey">
          <SingleLine placeholder="my-message-name" />
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

      <Tabs.TabPane forceRender tab={<TriggerTab hasContent={!!topic} label="Topic" />} key="topic">
        <Form.Item data-cy="topic" name="topic" rules={[{required: true, message: 'Please enter a topic'}]}>
          <SingleLine placeholder="my-topic" />
        </Form.Item>
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab={<TriggerTab totalItems={headers?.length} label="Headers" />} key="headers">
        <KeyValueList
          name="headers"
          label=""
          addButtonLabel="Add Header"
          keyPlaceholder="Header"
          valuePlaceholder="Value"
          initialValue={[{key: '', value: ''}]}
        />
      </Tabs.TabPane>

      <Tabs.TabPane
        forceRender
        tab={<TriggerTab hasContent={!!sslVerification || !!skipTraceCollection} label="Settings" />}
        key="settings"
      >
        <S.SettingsContainer>
          <SSL formID="kafka" />
          <SkipTraceCollection />
        </S.SettingsContainer>
      </Tabs.TabPane>
    </Tabs>
  );
};

export default Kafka;
