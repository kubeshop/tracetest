import { SupportedDataStores } from 'types/DataStore.types';
import AWSXRay from './ColorIcons/AwsXRay';
import AzureAppInsights from './ColorIcons/AzureAppInsights';
import Dash0 from './ColorIcons/Dash0';
import Datadog from './ColorIcons/Datadog';
import Dynatrace from './ColorIcons/Dynatrace';
import Elastic from './ColorIcons/Elastic';
import Honeycomb from './ColorIcons/Honeycomb';
import Instana from './ColorIcons/Instana';
import Jaeger from './ColorIcons/Jaeger';
import Lightstep from './ColorIcons/Lightstep';
import NewRelic from './ColorIcons/NewRelic';
import OpenSearch from './ColorIcons/OpenSearch';
import Otlp from './ColorIcons/Otlp';
import SignalFx from './ColorIcons/SignalFx';
import Signoz from './ColorIcons/Signoz';
import SumoLogic from './ColorIcons/SumoLogic';
import Tempo from './ColorIcons/Tempo';

const colorIconMap = {
  [SupportedDataStores.AWSXRay]: AWSXRay,
  [SupportedDataStores.AzureAppInsights]: AzureAppInsights,
  [SupportedDataStores.Dash0]: Dash0,
  [SupportedDataStores.Datadog]: Datadog,
  [SupportedDataStores.Dynatrace]: Dynatrace,
  [SupportedDataStores.ElasticApm]: Elastic,
  [SupportedDataStores.Honeycomb]: Honeycomb,
  [SupportedDataStores.Instana]: Instana,
  [SupportedDataStores.JAEGER]: Jaeger,
  [SupportedDataStores.Lightstep]: Lightstep,
  [SupportedDataStores.NewRelic]: NewRelic,
  [SupportedDataStores.OpenSearch]: OpenSearch,
  [SupportedDataStores.OtelCollector]: Otlp,
  [SupportedDataStores.SignalFX]: SignalFx,
  [SupportedDataStores.Signoz]: Signoz,
  [SupportedDataStores.SumoLogic]: SumoLogic,
  [SupportedDataStores.TEMPO]: Tempo,
} as const;

export default colorIconMap;
