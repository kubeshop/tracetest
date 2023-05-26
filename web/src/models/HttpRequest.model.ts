import {Model, THttpSchemas} from 'types/Common.types';
import {THeader} from 'types/Test.types';

export type TRawHTTPRequest = THttpSchemas['HTTPRequest'];
type HttpRequest = Model<
  TRawHTTPRequest,
  {
    headers: THeader[];
  }
>;

const HttpRequest = ({
  method = 'GET',
  url = '',
  headers = [],
  body = '',
  auth = {},
  sslVerification = false,
}: TRawHTTPRequest): HttpRequest => {
  return {
    method,
    url,
    headers: headers.map(({key = '', value = ''}) => ({
      key,
      value,
    })),
    body,
    auth,
    sslVerification,
  };
};

export default HttpRequest;
