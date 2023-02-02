import {Model, TTraceIDSchemas} from 'types/Common.types';

export type TRawTRACEIDRequest = TTraceIDSchemas['TRACEIDRequest'];
type TraceIDRequest = Model<
  TRawTRACEIDRequest,
  {
    id: string;
  }
>;

const TraceIDRequest = ({id = ''}: TRawTRACEIDRequest): TraceIDRequest => {
  return {id};
};

export default TraceIDRequest;
