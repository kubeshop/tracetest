import {THTTPRequest, TRawHTTPRequest} from '../types/Test.types';

const HttpRequest = ({method = 'GET', url = '', headers = [], body = '', auth = {}}: TRawHTTPRequest): THTTPRequest => {
  return {
    method,
    url,
    headers: headers.map(({key = '', value = ''}) => ({
      key,
      value,
    })),
    body,
    auth,
  };
};

export default HttpRequest;
