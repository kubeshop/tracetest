import {SupportedDataStores, TRawServerConfig, TServerConfig} from 'types/Config.types';

const ServerConfig = ({
  telemetry: {exporter = '', applicationExporter = '', dataStore = SupportedDataStores.JAEGER} = {
    exporter: '',
    applicationExporter: '',
    dataStore: SupportedDataStores.JAEGER,
  },
}: TRawServerConfig): TServerConfig => {
  return {
    telemetry: {
      exporter,
      applicationExporter,
      dataStore,
    },
  };
};

export default ServerConfig;
