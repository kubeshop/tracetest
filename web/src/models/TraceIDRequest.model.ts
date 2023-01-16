import {TRawTRACEIDRequest, TTRACEIDRequest} from '../types/Test.types';

const TraceIDRequest = ({id = ''}: TRawTRACEIDRequest): TTRACEIDRequest => {
  return {id};
};

export default TraceIDRequest;
