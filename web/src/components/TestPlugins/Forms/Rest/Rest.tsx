import {Form, Tabs} from 'antd';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import {Body} from 'components/Inputs';
import {Auth, SSL, KeyValueList, SkipTraceCollection} from 'components/Fields';
import * as S from './Rest.styled';

const Rest = () => (
  <Tabs defaultActiveKey="auth">
    <Tabs.TabPane forceRender tab="Auth" key="auth">
      <Auth />
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Body" key="body">
      <Form.Item name="body" noStyle>
        <Body />
      </Form.Item>
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Headers" key="headers">
      <KeyValueList
        name="headers"
        label=""
        addButtonLabel="Add Header"
        keyPlaceholder="Header"
        valuePlaceholder="Value"
        initialValue={DEFAULT_HEADERS}
      />
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Settings" key="settings">
      <S.SettingsContainer>
        <SSL formID="rest" />
        <SkipTraceCollection />
      </S.SettingsContainer>
    </Tabs.TabPane>
  </Tabs>
);

export default Rest;
