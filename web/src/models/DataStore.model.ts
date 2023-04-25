import {SupportedDataStores} from 'types/DataStore.types';
import {Model, TDataStoreSchemas} from 'types/Common.types';

export type TRawDataStore = TDataStoreSchemas['DataStoreResource'];
type DataStore = Model<TRawDataStore, {}>['spec'] & {otlp?: {}; newRelic?: {}; lightstep?: {}; datadog?: {}};

const DataStore = ({
  spec: {
    id = '',
    name = '',
    type = SupportedDataStores.JAEGER,
    default: isDefault = false,
    createdAt = '',
    openSearch = {},
    elasticApm = {},
    signalFx = {},
    jaeger = {},
    tempo = {},
    awsxray = {},
  } = {id: '', name: '', type: SupportedDataStores.JAEGER},
}: TRawDataStore): DataStore => ({
  id,
  name,
  type,
  default: isDefault,
  createdAt,
  openSearch,
  elasticApm,
  signalFx,
  jaeger,
  tempo,
  awsxray,
});

export default DataStore;
