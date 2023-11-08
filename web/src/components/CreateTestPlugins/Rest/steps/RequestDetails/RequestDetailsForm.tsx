import {Tabs} from 'antd';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import KeyValueListInput from 'components/KeyValueListInput';
import BodyField from './BodyField/BodyField';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/RequestDetailsAuthInput';
import SSLVerification from './SSLVerification';

export const FORM_ID = 'create-test';

const RequestDetailsForm = () => (
  <Tabs defaultActiveKey="auth">
    <Tabs.TabPane forceRender tab="Auth" key="auth">
      <RequestDetailsAuthInput />
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Body" key="body">
      <BodyField />
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
      <SSLVerification />
    </Tabs.TabPane>
  </Tabs>
);

export default RequestDetailsForm;
