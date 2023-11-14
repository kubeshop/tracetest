import {Form, Tabs} from 'antd';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import {Body} from 'components/Inputs';
import {Auth, KeyValueList, SSL} from 'components/Fields';
import {FORM_ID} from 'pages/CreateTest/Content';

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
      <SSL formID={FORM_ID} />
    </Tabs.TabPane>
  </Tabs>
);

export default RequestDetailsForm;
