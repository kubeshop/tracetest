import {TriggerTypes} from 'constants/Test.constants';
import TriggerHeaderBarGrpc from './Grpc';
import TriggerHeaderBarHttp from './Http';
import TriggerHeaderBarTraceID from './TraceID';
import TriggerHeaderBarKafka from './Kafka';

const TriggerHeaderMap = {
  [TriggerTypes.http]: TriggerHeaderBarHttp,
  [TriggerTypes.grpc]: TriggerHeaderBarGrpc,
  [TriggerTypes.traceid]: TriggerHeaderBarTraceID,
  [TriggerTypes.kafka]: TriggerHeaderBarKafka,
};

interface IProps {
  type: TriggerTypes;
}

const TriggerHeaderBar = ({type}: IProps) => {
  const Component = TriggerHeaderMap[type];

  return <Component />;
};

export default TriggerHeaderBar;
