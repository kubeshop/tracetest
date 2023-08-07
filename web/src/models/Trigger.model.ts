import {TriggerTypes} from 'constants/Test.constants';
import {get} from 'lodash';
import {TTriggerSchemas} from 'types/Common.types';
import {TTriggerRequest} from 'types/Test.types';
import GrpcRequest from './GrpcRequest.model';
import HttpRequest from './HttpRequest.model';
import TraceIDRequest from './TraceIDRequest.model';
import KafkaRequest from './KafkaRequest.model';

export type TRawTrigger = TTriggerSchemas['Trigger'];
type Trigger = {
  type: TriggerTypes;
  entryPoint: string;
  method: string;
  request: TTriggerRequest;
};

const EntryData = {
  [TriggerTypes.http](request: object) {
    return {
      entryPoint: get(request, 'url', ''),
      method: get(request, 'method', ''),
    };
  },
  [TriggerTypes.grpc](request: object) {
    return {
      entryPoint: get(request, 'address', ''),
      method: get(request, 'method', ''),
    };
  },
  [TriggerTypes.traceid](request: object) {
    return {
      entryPoint: get(request, 'id', ''),
      method: 'TraceID',
    };
  },
  [TriggerTypes.kafka](request: object) {
    return {
      entryPoint: get(request, 'brokerUrls', ''), // TODO (daniel)
      method: 'Kafka',
    };
  },
};

const Trigger = ({type: rawType = 'http', httpRequest = {}, grpc = {}, traceid = {}, kafka = {}}: TRawTrigger): Trigger => {
  const type = rawType as TriggerTypes;

  let request = {} as TTriggerRequest;
  if (type === TriggerTypes.http) {
    request = HttpRequest(httpRequest);
  } else if (type === TriggerTypes.grpc) {
    request = GrpcRequest(grpc);
  } else if (type === TriggerTypes.traceid) {
    request = TraceIDRequest(traceid);
  } else if (type === TriggerTypes.kafka) {
    request = KafkaRequest(kafka);
  }

  const {entryPoint, method} = EntryData[type || TriggerTypes.http](request);

  return {
    type,
    entryPoint,
    method,
    request,
  };
};

export default Trigger;
