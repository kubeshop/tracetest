import {Form, Tabs} from 'antd';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import {Body} from 'components/Inputs';
import TriggerTab from 'components/TriggerTab';
import useQueryTabs from 'hooks/useQueryTabs';
import {Auth, SSL, KeyValueList, SkipTraceCollection} from 'components/Fields';
import {TDraftTest} from 'types/Test.types';
import * as S from './Rest.styled';

const Rest = () => {
  const [activeKey, setActiveKey] = useQueryTabs('auth', 'triggerTab');
  const form = Form.useFormInstance<TDraftTest>();
  const authType = Form.useWatch(['auth', 'type'], form);
  const body = Form.useWatch('body', form);
  const headers = Form.useWatch('headers', form);
  const sslVerification = Form.useWatch('sslVerification', form);
  const skipTraceCollection = Form.useWatch('skipTraceCollection', form);

  return (
    <Tabs defaultActiveKey={activeKey} activeKey={activeKey} onChange={setActiveKey}>
      <Tabs.TabPane forceRender tab={<TriggerTab hasContent={!!authType} label="Auth" />} key="auth">
        <Auth />
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab={<TriggerTab hasContent={!!body} label="Body" />} key="body">
        <Form.Item name="body" noStyle>
          <Body />
        </Form.Item>
      </Tabs.TabPane>

      <Tabs.TabPane forceRender tab={<TriggerTab totalItems={headers?.length} label="Headers" />} key="headers">
        <KeyValueList
          name="headers"
          label=""
          addButtonLabel="Add Header"
          keyPlaceholder="Header"
          valuePlaceholder="Value"
          initialValue={DEFAULT_HEADERS}
        />
      </Tabs.TabPane>

      <Tabs.TabPane
        forceRender
        tab={<TriggerTab hasContent={!!sslVerification || !!skipTraceCollection} label="Settings" />}
        key="settings"
      >
        <S.SettingsContainer>
          <SSL formID="rest" />
          <SkipTraceCollection />
        </S.SettingsContainer>
      </Tabs.TabPane>
    </Tabs>
  );
};

export default Rest;
