import {useTheme} from 'styled-components';
import {SupportedDataStores} from 'types/DataStore.types';
import Agent from './Icons/Agent';
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

const iconMap = {
  [SupportedDataStores.Agent]: Agent,
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
} as const;

interface IProps {
  color?: string;
  dataStoreType: SupportedDataStores;
  width?: string;
  height?: string;
}

export interface IIconProps {
  color: string;
  width?: string;
  height?: string;
}

const DataStoreIcon = ({color, dataStoreType, width, height}: IProps) => {
  const {
    color: {text: defaultColor},
  } = useTheme();
  const Component = iconMap[dataStoreType];

  return <Component color={color || defaultColor} width={width} height={height} />;
};

export default DataStoreIcon;
