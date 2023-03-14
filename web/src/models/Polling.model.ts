import {Model, TConfigSchemas} from 'types/Common.types';

export type TRawPolling = TConfigSchemas['PollingProfile'];
type Polling = Model<Model<TRawPolling, {}>['spec'], {}>;

function Polling({
  spec: {
    default: isDefault = false,
    id = '',
    name = '',
    periodic: {retryDelay = '', timeout = ''} = {},
    strategy = 'periodic',
  } = {
    id: '',
    name: '',
    strategy: 'periodic',
  },
}: TRawPolling = {}): Polling {
  return {
    default: isDefault,
    id,
    name,
    periodic: {
      retryDelay,
      timeout,
    },
    strategy,
  };
}

export default Polling;
