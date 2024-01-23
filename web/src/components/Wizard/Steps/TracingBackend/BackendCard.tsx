import DataStoreIcon from 'components/DataStoreIcon/DataStoreIcon';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import {SupportedDataStores} from 'types/DataStore.types';
import * as S from './TracingBackend.styled';

interface IProps {
  backend: SupportedDataStores;
  onSelect(): void;
  selectedBackend?: SupportedDataStores;
}

const BackendCard = ({backend, onSelect, selectedBackend}: IProps) => (
  <S.Card onClick={onSelect}>
    <DataStoreIcon withColor dataStoreType={backend} height="22" width="22" />{' '}
    <S.Name style={{textOverflow: 'ellipsis'}}>{SupportedDataStoresToName[backend]}</S.Name>
    {selectedBackend === backend && <S.SelectedIcon />}
  </S.Card>
);

export default BackendCard;
