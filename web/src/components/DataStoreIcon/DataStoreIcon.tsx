import {useTheme} from 'styled-components';
import {SupportedDataStores} from 'types/Config.types';
import Elastic from './Icons/Elastic';
import Jaeger from './Icons/Jaeger';
import Lightstep from './Icons/Lightstep';
import Datadog from './Icons/Datadog';
import NewRelic from './Icons/NewRelic';
import OpenSearch from './Icons/OpenSearch';
import Otlp from './Icons/Otlp';
import SignalFx from './Icons/SignalFx';
import Tempo from './Icons/Tempo';

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
  [SupportedDataStores.AWSXRay]: Datadog,
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
