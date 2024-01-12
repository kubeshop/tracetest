import {SupportedDataStores} from 'types/DataStore.types';
import Elastic from './Icons/Elastic';
import Jaeger from './Icons/Jaeger';
import Lightstep from './Icons/Lightstep';
import Datadog from './Icons/Datadog';
import NewRelic from './Icons/NewRelic';
import OpenSearch from './Icons/OpenSearch';
import Otlp from './Icons/Otlp';
import SignalFx from './Icons/SignalFx';
import Tempo from './Icons/Tempo';
import AWSXRay from './Icons/AwsXRay';
import Honeycomb from './Icons/Honeycomb';
import AzureAppInsights from './Icons/AzureAppInsights';
import Signoz from './Icons/Signoz';
import Dynatrace from './Icons/Dynatrace';
import SumoLogic from './Icons/SumoLogic';

const iconMap = {
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

export default iconMap;