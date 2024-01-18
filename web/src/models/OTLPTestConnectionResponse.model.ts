import {Model, TConfigSchemas} from 'types/Common.types';

export type TRawOTLPTestConnectionResponse = TConfigSchemas['OTLPTestConnectionResponse'];
type OTLPTestConnectionResponse = Model<TRawOTLPTestConnectionResponse, {}>;

const defaultOTLPTestConnectionResponse: TRawOTLPTestConnectionResponse = {
  spanCount: 0,
  lastSpanTimestamp: '',
};

function OTLPTestConnectionResponse({
  spanCount = 0,
  lastSpanTimestamp = '',
} = defaultOTLPTestConnectionResponse): OTLPTestConnectionResponse {
  return {
    spanCount,
    lastSpanTimestamp,
  };
}

export default OTLPTestConnectionResponse;
