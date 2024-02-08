import {SupportedDataStores} from 'types/DataStore.types';
import AWSXRay from './Icons/AwsXRay';
import AzureAppInsights from './Icons/AzureAppInsights';
import Datadog from './Icons/Datadog';
import Dynatrace from './Icons/Dynatrace';
import Elastic from './Icons/Elastic';
import Honeycomb from './Icons/Honeycomb';
import Instana from './Icons/Instana';
import Jaeger from './Icons/Jaeger';
import Lightstep from './Icons/Lightstep';
import NewRelic from './Icons/NewRelic';
import OpenSearch from './Icons/OpenSearch';
import Otlp from './Icons/Otlp';
import SignalFx from './Icons/SignalFx';
import Signoz from './Icons/Signoz';
import SumoLogic from './Icons/SumoLogic';
import Tempo from './Icons/Tempo';

const iconMap = {
  [SupportedDataStores.AWSXRay]: AWSXRay,
  [SupportedDataStores.AzureAppInsights]: AzureAppInsights,
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

export default iconMap;
