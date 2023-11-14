import {Form, Tabs} from 'antd';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import KeyValueListInput from 'components/Fields/KeyValueList';
import {Body} from 'components/Inputs';
import {Auth, SSL, SkipTraceCollection} from 'components/Fields';
import * as S from './RequestDetails.styled';

export const FORM_ID = 'create-test';

const RequestDetailsForm = () => (
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
      <KeyValueListInput
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
        <SSL />
        <SkipTraceCollection />
      </S.SettingsContainer>
    </Tabs.TabPane>
  </Tabs>
);

export default RequestDetailsForm;
