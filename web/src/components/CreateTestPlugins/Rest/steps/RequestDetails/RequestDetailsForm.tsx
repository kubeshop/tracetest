import {Form, Tabs} from 'antd';
import {IHttpValues, TDraftTestForm} from 'types/Test.types';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import KeyValueListInput from 'components/KeyValueListInput';
import BodyField from './BodyField/BodyField';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/RequestDetailsAuthInput';
import SSLVerification from './SSLVerification';

export const FORM_ID = 'create-test';

interface IProps {
  form: TDraftTestForm<IHttpValues>;
}

const RequestDetailsForm = ({form}: IProps) => (
  <Tabs defaultActiveKey="general">
    <Tabs.TabPane forceRender tab="General" key="general">
      <BodyField setBody={body => form.setFieldsValue({body})} body={Form.useWatch('body', form)} />
      <RequestDetailsAuthInput />
      <SSLVerification />
    </Tabs.TabPane>

    <Tabs.TabPane forceRender tab="Headers" key="headers">
      <KeyValueListInput
        name="headers"
        label="Header list"
        addButtonLabel="Add Header"
        keyPlaceholder="Header"
        valuePlaceholder="Value"
        initialValue={DEFAULT_HEADERS}
      />
    </Tabs.TabPane>
  </Tabs>
);

export default RequestDetailsForm;
