import {SupportedDataStores} from 'types/Config.types';
import DataStorePlugin from '.';

interface IProps {
  dataStoreType?: SupportedDataStores;
}

const DataStoreComponentFactory = ({dataStoreType = SupportedDataStores.JAEGER}: IProps) => {
  const FormComponent = DataStorePlugin[dataStoreType];

  return <FormComponent />;
};

export default DataStoreComponentFactory;
