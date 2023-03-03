import {Popover} from 'antd';
import {noop} from 'lodash';
import {useTheme} from 'styled-components';
import {ConfigMode, SupportedDataStores} from 'types/DataStore.types';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import {useDataStoreConfig} from 'providers/DataStoreConfig/DataStoreConfig.provider';

import DataStoreIcon from '../../DataStoreIcon/DataStoreIcon';
import * as S from './DataStoreForm.styled';

interface IProps {
  value?: SupportedDataStores;
  onChange?(value: SupportedDataStores): void;
}

const supportedDataStoreList = Object.values(SupportedDataStores);

const DataStoreSelectionInput = ({onChange = noop, value = SupportedDataStores.JAEGER}: IProps) => {
  const {
    color: {text, primary},
  } = useTheme();
  const {dataStoreConfig} = useDataStoreConfig();
  const configuredDataStoreType = dataStoreConfig.defaultDataStore.type;

  return (
    <S.DataStoreListContainer>
      {supportedDataStoreList.map(dataStore => {
        const isSelected = value === dataStore;
        const isConfigured = configuredDataStoreType === dataStore && dataStoreConfig.mode === ConfigMode.READY;
        return (
          <S.DataStoreItemContainer $isSelected={isSelected} key={dataStore} onClick={() => onChange(dataStore)}>
            <DataStoreIcon dataStoreType={dataStore} color={isSelected ? primary : text} width="22" height="22" />
            <S.DataStoreName $isSelected={isSelected}>{SupportedDataStoresToName[dataStore]}</S.DataStoreName>
            {isConfigured && (
              <Popover content="This data source is currently configured" placement="right">
                <S.InfoIcon />
              </Popover>
            )}
          </S.DataStoreItemContainer>
        );
      })}
    </S.DataStoreListContainer>
  );
};

export default DataStoreSelectionInput;
