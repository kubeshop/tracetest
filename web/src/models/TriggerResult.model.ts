import {TriggerTypes} from 'constants/Test.constants';
import {get} from 'lodash';
import {THeader} from 'types/Test.types';
import {TTriggerSchemas} from 'types/Common.types';

export type TRawTriggerResult = TTriggerSchemas['TriggerResult'];
type TriggerResult = {
  type: TriggerTypes;
  headers?: THeader[];
  body?: string;
  statusCode: number;
  bodyMimeType?: string;
};

const ResponseData = {
  [TriggerTypes.http](response: object) {
    return {
      body: get(response, 'body', ''),
      headers: get(response, 'headers', undefined) as THeader[],
      statusCode: get(response, 'statusCode', 200),
    };
  },
  [TriggerTypes.grpc](response: object) {
    return {
      body: get(response, 'body', ''),
      headers: get(response, 'metadata', undefined) as THeader[],
      statusCode: get(response, 'statusCode', 0),
    };
  },
  [TriggerTypes.traceid](response: object) {
    return {
      body: get(response, 'id', ''),
      headers: [],
      statusCode: 200,
    };
  },
};

const TriggerResult = ({
  type: rawType = 'http',
  triggerResult: {http = {}, grpc = {statusCode: 0}, traceid = {}} = {},
}: TRawTriggerResult): TriggerResult => {
  const type = rawType as TriggerTypes;

  let request = {};
  if (type === TriggerTypes.http) {
    request = http;
  } else if (type === TriggerTypes.grpc) {
    request = grpc;
  } else if (type === TriggerTypes.traceid) {
    request = traceid;
  }

  const {body, headers = [], statusCode} = ResponseData[type](request);

  const bodyMimeType = headers.find(({key}) => key.toLowerCase() === 'content-type')?.value;

  return {
    type,
    body,
    headers,
    statusCode,
    bodyMimeType,
  };
};

export default TriggerResult;
