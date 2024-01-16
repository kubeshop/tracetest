import { SupportedDataStores } from 'types/DataStore.types';
import Elastic from './ColorIcons/Elastic';
import Jaeger from './ColorIcons/Jaeger';
import Lightstep from './ColorIcons/Lightstep';
import Datadog from './ColorIcons/Datadog';
import NewRelic from './ColorIcons/NewRelic';
import OpenSearch from './ColorIcons/OpenSearch';
import Otlp from './ColorIcons/Otlp';
import SignalFx from './ColorIcons/SignalFx';
import Tempo from './ColorIcons/Tempo';
import AWSXRay from './ColorIcons/AwsXRay';
import Honeycomb from './ColorIcons/Honeycomb';
import AzureAppInsights from './ColorIcons/AzureAppInsights';
import Signoz from './ColorIcons/Signoz';
import Dynatrace from './ColorIcons/Dynatrace';
import SumoLogic from './ColorIcons/SumoLogic';

const colorIconMap = {
  [SupportedDataStores.JAEGER]: Jaeger,
  [SupportedDataStores.SignalFX]: SignalFx,
  [SupportedDataStores.ElasticApm]: Elastic,
  [SupportedDataStores.OtelCollector]: Otlp,
  [SupportedDataStores.TEMPO]: Tempo,
  [SupportedDataStores.OpenSearch]: OpenSearch,
  [SupportedDataStores.NewRelic]: NewRelic,
  [SupportedDataStores.Lightstep]: Lightstep,
  [SupportedDataStores.Datadog]: Datadog,
  [SupportedDataStores.AWSXRay]: AWSXRay,
  [SupportedDataStores.Honeycomb]: Honeycomb,
  [SupportedDataStores.AzureAppInsights]: AzureAppInsights,
  [SupportedDataStores.Signoz]: Signoz,
  [SupportedDataStores.Dynatrace]: Dynatrace,
  [SupportedDataStores.SumoLogic]: SumoLogic,
} as const;

export default colorIconMap;
