import {TriggerTypes} from 'constants/Test.constants';
import {get} from 'lodash';
import {TTriggerSchemas} from 'types/Common.types';
import {TTriggerRequest} from 'types/Test.types';
import GrpcRequest from './GrpcRequest.model';
import HttpRequest from './HttpRequest.model';
import TraceIDRequest from './TraceIDRequest.model';

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
};

const Trigger = ({triggerType = 'http', http = {}, grpc = {}, traceid = {}}: TRawTrigger): Trigger => {
  const type = triggerType as TriggerTypes;

  let request = {} as TTriggerRequest;
  if (type === TriggerTypes.http) {
    request = HttpRequest(http);
  } else if (type === TriggerTypes.grpc) {
    request = GrpcRequest(grpc);
  } else if (type === TriggerTypes.traceid) {
    request = TraceIDRequest(traceid);
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
