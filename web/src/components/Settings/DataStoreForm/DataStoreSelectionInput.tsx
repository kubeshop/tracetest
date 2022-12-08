import {capitalize, noop} from 'lodash';
import {SupportedDataStores} from 'types/Config.types';
import * as S from './DataStoreForm.styled';

interface IProps {
  value?: SupportedDataStores;
  onChange?(value: SupportedDataStores): void;
}

const supportedDataStoreList = Object.values(SupportedDataStores);

const DataStoreSelectionInput = ({onChange = noop, value = SupportedDataStores.JAEGER}: IProps) => {
  return (
    <S.DataStoreListContainer>
      {supportedDataStoreList.map(dataStore => {
        const isSelected = value === dataStore;
        return (
          <S.DataStoreItemContainer $isSelected={isSelected} key={dataStore} onClick={() => onChange(dataStore)}>
            <S.Circle>{isSelected && <S.Check />}</S.Circle>
            <S.DataStoreIcon $dataStore={dataStore} />
            <S.DataStoreName>{capitalize(dataStore)}</S.DataStoreName>
          </S.DataStoreItemContainer>
        );
      })}
    </S.DataStoreListContainer>
  );
};

export default DataStoreSelectionInput;
