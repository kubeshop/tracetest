import BasicDetailsForm from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsForm';
import {IHttpValues, TDraftTestForm} from 'types/Test.types';
import {IFormProps} from '../EditRequestDetails';

const EditRequestDetailsHttp = ({form}: IFormProps) => <BasicDetailsForm form={form as TDraftTestForm<IHttpValues>} />;

export default EditRequestDetailsHttp;
