import {SupportedDataStores} from 'types/DataStore.types';
import {Model, TDataStoreSchemas} from 'types/Common.types';

export type TRawDataStore = TDataStoreSchemas['DataStoreResource'];
type DataStore = Model<TRawDataStore, {}>['spec'] & {otlp?: {}; newrelic?: {}; lightstep?: {}; datadog?: {}};

const DataStore = ({
  spec: {
    id = '',
    name = '',
    type = SupportedDataStores.JAEGER,
    default: isDefault = false,
    createdAt = '',
    opensearch = {},
    elasticapm = {},
    signalfx = {},
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
  opensearch,
  elasticapm,
  signalfx,
  jaeger,
  tempo,
  awsxray,
});

export default DataStore;
