import {TRawTelemetryConfig, TTelemetryConfig} from 'types/Config.types';

const TelemetryConfig = ({exporters = [], dataStores = []}: TRawTelemetryConfig): TTelemetryConfig => {
  // todo add data stores and exporters models

  return {
    exporters,
    dataStores,
  };
};

export default TelemetryConfig;
