interface IConfigCommon {
  dataStoreName: string;
  dataStoreType: string;
  poolingMaxWaitTimeForTrace: string;
  poolingRetryDelay: string;
}

interface IDataStoreGRPC {
  // TODO
}

interface IDataStoreOpenSearch {
  // TODO
}

interface IDataStoreSignalFx {
  // TODO
}

type TDataStoreUnion = IDataStoreGRPC | IDataStoreOpenSearch | IDataStoreSignalFx;

export type TConfig<T = TDataStoreUnion> = Partial<IConfigCommon & T>;

// DataStore service interface
export interface IDataStoreService {
  getYamlConfig(values: TConfig): void; // TODO: define return type
  validate(values: TConfig): boolean;
}
