import {DataStores} from 'constants/Config.constants';
import {TConfig} from 'types/Config.types';
import GRPCService from './DataStores/GRPC.service';

const DataStoreServices = {
  [DataStores.GRPC]: GRPCService,
  [DataStores.OPEN_SEARCH]: GRPCService, // TODO
  [DataStores.SIGNAL_FX]: GRPCService, // TODO
} as const;

const ConfigService = {
  getYamlConfig(values: TConfig) {
    // TODO
  },
  validate(values: TConfig) {
    // TODO
  },
};

export default ConfigService;
