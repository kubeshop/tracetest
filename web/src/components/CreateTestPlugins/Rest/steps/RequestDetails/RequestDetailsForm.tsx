import {Tabs} from 'antd';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import KeyValueListInput from 'components/Fields/KeyValueList';
import {Body} from 'components/Inputs';
import {Auth, SSL} from 'components/Fields';

export const FORM_ID = 'create-test';

const RequestDetailsForm = () => (
  <Tabs defaultActiveKey="auth">
    <Tabs.TabPane forceRender tab="Auth" key="auth">
      <Auth />
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Body" key="body">
      <Body />
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
      <SSL />
    </Tabs.TabPane>
  </Tabs>
);

export default RequestDetailsForm;
