import {TriggerTypes} from 'constants/Test.constants';
import {TDraftTestForm, TTriggerRequest} from 'types/Test.types';
import EditRequestDetailsHttp from './Http';
import EditRequestDetailsGrpc from './Grpc';

const EditRequestDetailsMap = {
  [TriggerTypes.http]: EditRequestDetailsHttp,
  [TriggerTypes.grpc]: EditRequestDetailsGrpc,
};

export interface IFormProps {
  form: TDraftTestForm;
  request: TTriggerRequest;
}

interface IProps {
  type: TriggerTypes;
  form: TDraftTestForm;
  request: TTriggerRequest;
}

const EditRequestDetails = ({type, form, request}: IProps) => {
  const Component = EditRequestDetailsMap[type];

  return <Component form={form} request={request} />;
};

export default EditRequestDetails;
