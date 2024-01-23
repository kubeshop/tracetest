import {useTheme} from 'styled-components';
import {SupportedDataStores} from 'types/DataStore.types';
import colorIconMap from './colorIconMap';
import iconMap from './iconMap';

interface IProps {
  color?: string;
  dataStoreType: SupportedDataStores;
  width?: string;
  height?: string;
  withColor?: boolean;
}

export interface IIconProps {
  color: string;
  width?: string;
  height?: string;
}

const DataStoreIcon = ({color, withColor = false, dataStoreType, width, height}: IProps) => {
  const {
    color: {text: defaultColor},
  } = useTheme();
  const Component = withColor ? colorIconMap[dataStoreType] : iconMap[dataStoreType];

  return <Component color={color || defaultColor} width={width} height={height} />;
};

export default DataStoreIcon;
