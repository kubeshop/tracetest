import {TriggerTypes} from 'constants/Test.constants';
import TriggerHeaderBarGrpc from './EntryPoint/Grpc';
import TriggerHeaderBarHttp from './EntryPoint/Http';
import TriggerHeaderBarTraceID from './EntryPoint/TraceID';
import TriggerHeaderBarKafka from './EntryPoint/Kafka';

const EntryPointFactoryMap = {
  [TriggerTypes.http]: TriggerHeaderBarHttp,
  [TriggerTypes.grpc]: TriggerHeaderBarGrpc,
  [TriggerTypes.kafka]: TriggerHeaderBarKafka,
  [TriggerTypes.traceid]: TriggerHeaderBarTraceID,
  [TriggerTypes.cypress]: () => null,
  [TriggerTypes.playwright]: () => null,
};

interface IProps {
  type: TriggerTypes;
}

const EntryPointFactory = ({type}: IProps) => {
  const Component = EntryPointFactoryMap[type];

  return <Component />;
};

export default EntryPointFactory;
