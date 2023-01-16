import {TriggerTypes} from 'constants/Test.constants';
import {get} from 'lodash';
import {TRawTrigger, TTrigger} from 'types/Test.types';
import GrpcRequest from './GrpcRequest.model';
import HttpRequest from './HttpRequest.model';
import TraceIDRequest from './TraceIDRequest.model';

const entryPointMap = {
  [TriggerTypes.http]: 'url',
  [TriggerTypes.grpc]: 'address',
  [TriggerTypes.traceid]: 'id',
} as const;

const entryMethodMap = {
  [TriggerTypes.http]: 'method',
  [TriggerTypes.grpc]: 'method',
  [TriggerTypes.traceid]: 'id',
} as const;

const getEntryData = (type: TriggerTypes, request: object) => {
  const entryPointField = entryPointMap[type];
  const entryMethodField = entryMethodMap[type];

  return {
    entryPoint: get(request, entryPointField, ''),
    method: get(request, entryMethodField, ''),
  };
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

  const {entryPoint, method} = getEntryData(type, request);

  return {
    type,
    entryPoint,
    method,
    request,
  };
};

export default Trigger;
