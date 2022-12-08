import {IMockFactory} from 'types/Common.types';
import {TConfig, TRawConfig} from 'types/Config.types';
import Config from '../Config.model';

const ConfigMock: IMockFactory<TConfig, TRawConfig> = () => ({
  raw(data = {}) {
    return {
      telemetry: {
        dataStores: [],
        exporters: [],
      },
      server: {
        telemetry: {
          dataStore: 'jaeger',
          exporter: '',
          applicationExporter: '',
        },
      },
      ...data,
    };
  },
  model(data = {}) {
    return Config(this.raw(data));
  },
});

export default ConfigMock();
