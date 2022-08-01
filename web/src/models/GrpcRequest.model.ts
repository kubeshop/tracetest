import {TGRPCRequest, TRawGRPCRequest} from '../types/Test.types';

const GrpcRequest = ({
  protobufFile = '',
  address = '',
  service = '',
  method = '',
  metadata = [],
  auth = {},
  request = '',
}: TRawGRPCRequest): TGRPCRequest => {
  return {
    protobufFile,
    address,
    service,
    method,
    metadata: metadata.map(({key = '', value = ''}) => ({
      key,
      value,
    })),
    auth,
    request,
  };
};

export default GrpcRequest;
