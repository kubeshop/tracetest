import {Model, TGrpcSchemas} from '../types/Common.types';
import {TRawGRPCHeader} from '../types/Test.types';

export type TRawGRPCRequest = TGrpcSchemas['GRPCRequest'];
type GrpcRequest = Model<
  TRawGRPCRequest,
  {
    metadata: Model<TRawGRPCHeader, {}>[];
  }
>;

const GrpcRequest = ({
  protobufFile = '',
  address = '',
  service = '',
  method = '',
  metadata = [],
  auth = {},
  request = '',
}: TRawGRPCRequest): GrpcRequest => {
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
