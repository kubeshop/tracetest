import {SupportedDataStoresToDocsLink, SupportedDataStoresToName} from 'constants/DataStore.constants';
import {SupportedDataStores} from 'types/DataStore.types';
import DocsBanner from 'components/DocsBanner/DocsBanner';

interface IProps {
  dataStoreType: SupportedDataStores;
}

const DataStoreDocsBanner = ({dataStoreType}: IProps) => {
  return (
    <DocsBanner>
      Need more information about setting up {SupportedDataStoresToName[dataStoreType]}?{' '}
      <a href={SupportedDataStoresToDocsLink[dataStoreType]} target="_blank">
        Go to our docs
      </a>
    </DocsBanner>
  );
};

export default DataStoreDocsBanner;
