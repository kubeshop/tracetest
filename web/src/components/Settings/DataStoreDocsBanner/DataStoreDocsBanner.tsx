import {ReadOutlined} from '@ant-design/icons';
import {SupportedDataStoresToDocsLink, SupportedDataStoresToName} from 'constants/DataStore.constants';
import {SupportedDataStores} from 'types/DataStore.types';
import * as S from './DataStoreDocsBanner.styled';

interface IProps {
  dataStoreType: SupportedDataStores;
}

const DataStoreDocsBanner = ({dataStoreType}: IProps) => {
  return (
    <S.DataStoreDocsBannerContainer>
      <ReadOutlined />
      <S.Text>
        Need more information about setting up {SupportedDataStoresToName[dataStoreType]}?{' '}
        <a href={SupportedDataStoresToDocsLink[dataStoreType]} target="_blank">
          Go to our docs
        </a>
      </S.Text>
    </S.DataStoreDocsBannerContainer>
  );
};

export default DataStoreDocsBanner;
