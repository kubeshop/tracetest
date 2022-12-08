import {TRawServerConfig, TServerConfig} from 'types/Config.types';

const ServerConfig = ({
  telemetry: {exporter = '', applicationExporter = '', dataStore = ''} = {
    exporter: '',
    applicationExporter: '',
    dataStore: '',
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
