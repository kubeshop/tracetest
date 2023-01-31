import {TriggerTypes} from 'constants/Test.constants';
import {get} from 'lodash';
import {TRawTrigger, TTrigger} from 'types/Test.types';
import GrpcRequest from './GrpcRequest.model';
import HttpRequest from './HttpRequest.model';
import TraceIDRequest from './TraceIDRequest.model';

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

const Trigger = ({
  triggerType = 'http',
  triggerSettings: {http = {}, grpc = {}, traceid = {}} = {},
}: TRawTrigger): TTrigger => {
  const type = triggerType as TriggerTypes;

  let request = {};
  if (type === TriggerTypes.http) {
    request = HttpRequest(http);
  } else if (type === TriggerTypes.grpc) {
    request = GrpcRequest(grpc);
  } else if (type === TriggerTypes.traceid) {
    request = TraceIDRequest(traceid);
  }

  const {entryPoint, method} = EntryData[type](request);

  return {
    type,
    entryPoint,
    method,
    request,
  };
};

export default Trigger;
