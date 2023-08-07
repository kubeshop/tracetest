import {TriggerTypes} from 'constants/Test.constants';
import {TDraftTestForm, TTriggerRequest} from 'types/Test.types';
import EditRequestDetailsHttp from './Http';
import EditRequestDetailsGrpc from './Grpc';
import EditDetailsTraceID from './TraceID';
import EditDetailsKafka from './Kafka';

const EditRequestDetailsMap = {
  [TriggerTypes.http]: EditRequestDetailsHttp,
  [TriggerTypes.grpc]: EditRequestDetailsGrpc,
  [TriggerTypes.traceid]: EditDetailsTraceID,
  [TriggerTypes.kafka]: EditDetailsKafka,
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
