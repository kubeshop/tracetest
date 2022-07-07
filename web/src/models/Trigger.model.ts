import {TriggerTypes} from 'constants/Test.constants';
import {get} from 'lodash';
import {TRawTrigger, TTrigger} from 'types/Test.types';
import GrpcRequest from './GrpcRequest.model';
import HttpRequest from './HttpRequest.model';

const entryPointMap = {
  [TriggerTypes.http]: 'url',
  [TriggerTypes.grpc]: 'address',
} as const;

const entryMethodMap = {
  [TriggerTypes.http]: 'method',
  [TriggerTypes.grpc]: 'method',
} as const;

const getEntryData = (type: TriggerTypes, request: object) => {
  const entryPointField = entryPointMap[type];
  const entryMethodField = entryMethodMap[type];

  return {
    entryPoint: get(request, entryPointField, ''),
    method: get(request, entryMethodField, ''),
  };
};

const Trigger = ({triggerType = 'http', triggerSettings: {http = {}, grpc = {}} = {}}: TRawTrigger): TTrigger => {
  const type = triggerType as TriggerTypes;
  const request = type === TriggerTypes.http ? HttpRequest(http) : GrpcRequest(grpc);
  const {entryPoint, method} = getEntryData(type, request);

  return {
    type,
    entryPoint,
    method,
    request,
  };
};

export default Trigger;
