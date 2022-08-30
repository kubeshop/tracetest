import {TriggerTypes} from 'constants/Test.constants';
import {get} from 'lodash';
import {TRawTriggerResult, TTriggerResult} from 'types/Test.types';

const entryBodyMap = {
  [TriggerTypes.http]: 'body',
  [TriggerTypes.grpc]: 'body',
} as const;

const entryHeadersMap = {
  [TriggerTypes.http]: 'headers',
  [TriggerTypes.grpc]: 'metadata',
} as const;

const entryStatusCodeMap = {
  [TriggerTypes.http]: 'statusCode',
  [TriggerTypes.grpc]: 'statusCode',
} as const;

const getResponseData = (type: TriggerTypes, response: object) => {
  const entryBodyField = entryBodyMap[type];
  const entryHeadersField = entryHeadersMap[type];
  const entryStatusCodeField = entryStatusCodeMap[type];

  return {
    body: get(response, entryBodyField, ''),
    headers: get(response, entryHeadersField, undefined),
    statusCode: get(response, entryStatusCodeField, 200),
  };
};

const TriggerResult = ({
  triggerType = 'http',
  triggerResult: {http = {}, grpc = {}} = {},
}: TRawTriggerResult): TTriggerResult => {
  const type = triggerType as TriggerTypes;
  const request = type === TriggerTypes.http ? http : grpc;
  const {body, headers, statusCode} = getResponseData(type, request);

  return {
    type,
    body,
    headers,
    statusCode,
  };
};

export default TriggerResult;
