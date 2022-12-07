import {ConfigMode, TConfig, TRawConfig} from 'types/Config.types';
import ServerConfig from './ServerConfig.model';
import TelemetryConfig from './TelemetryConfig.model';

const Config = ({server: rawServer = {}, telemetry: rawTelemetry = {}}: TRawConfig): TConfig => {
  const server = ServerConfig(rawServer);
  const telemetry = TelemetryConfig(rawTelemetry);
  const dataStoreType = server.telemetry.dataStore;
  const mode =
    (Boolean(dataStoreType && telemetry.dataStores.find(({type}) => type === dataStoreType)) && ConfigMode.READY) ||
    ConfigMode.NO_TRACING_MODE;

  return {
    server,
    telemetry,
    mode,
  };
};

export default Config;
