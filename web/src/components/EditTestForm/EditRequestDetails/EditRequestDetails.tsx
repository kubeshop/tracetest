import {TriggerTypes} from 'constants/Test.constants';
import {TDraftTestForm} from 'types/Test.types';
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
}

interface IProps {
  type: TriggerTypes;
  form: TDraftTestForm;
}

const EditRequestDetails = ({type, form}: IProps) => {
  const Component = EditRequestDetailsMap[type];

  return <Component form={form} />;
};

export default EditRequestDetails;
