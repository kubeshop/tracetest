import {TriggerTypes} from 'constants/Test.constants';
import {TDraftTestForm} from 'types/Test.types';
import TriggerHeaderBarGrpc from './Grpc';
import TriggerHeaderBarHttp from './Http';
import TriggerHeaderBarTraceID from './TraceID';

const TriggerHeaderMap = {
  [TriggerTypes.http]: TriggerHeaderBarHttp,
  [TriggerTypes.grpc]: TriggerHeaderBarGrpc,
  [TriggerTypes.traceid]: TriggerHeaderBarTraceID,
  [TriggerTypes.kafka]: TriggerHeaderBarHttp,
};

export interface IFormProps {
  form: TDraftTestForm;
}

interface IProps {
  form: TDraftTestForm;
  type: TriggerTypes;
}

const TriggerHeaderBar = ({form, type}: IProps) => {
  const Component = TriggerHeaderMap[type];

  return <Component form={form} />;
};

export default TriggerHeaderBar;
