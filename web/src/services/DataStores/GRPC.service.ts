import {IDataStoreService, TConfig} from 'types/Config.types';

const GRPCService: IDataStoreService = {
  getYamlConfig(values: TConfig) {
    // TODO
  },
  validate(values: TConfig): boolean {
    // TODO
    return true;
  },
};

export default GRPCService;
