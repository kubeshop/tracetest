import {Popover, Tabs} from 'antd';
import {noop} from 'lodash';
import {useTheme} from 'styled-components';
import {ConfigMode, SupportedDataStores} from 'types/DataStore.types';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';

import DataStoreIcon from '../../DataStoreIcon/DataStoreIcon';
import * as S from './DataStoreSelection.styled';

interface IProps {
  value?: SupportedDataStores;
  onChange?(value: SupportedDataStores): void;
}

const supportedDataStoreList = Object.values(SupportedDataStores);

const DataStoreSelection = ({onChange = noop, value = SupportedDataStores.JAEGER}: IProps) => {
  const {
    color: {text, primary},
  } = useTheme();
  const {dataStoreConfig} = useSettingsValues();
  const configuredDataStoreType = dataStoreConfig.defaultDataStore.type;

  return (
    <S.DataStoreListContainer
      tabPosition="left"
      onChange={dataStore => onChange(dataStore as SupportedDataStores)}
      defaultActiveKey={value}
    >
      {supportedDataStoreList.map(dataStore => {
        const isSelected = value === dataStore;
        const isConfigured = configuredDataStoreType === dataStore && dataStoreConfig.mode === ConfigMode.READY;

        return (
          <Tabs.TabPane
            key={dataStore}
            tab={
              <S.DataStoreItemContainer $isSelected={isSelected} key={dataStore}>
                <DataStoreIcon dataStoreType={dataStore} color={isSelected ? primary : text} width="22" height="22" />
                <S.DataStoreName $isSelected={isSelected}>{SupportedDataStoresToName[dataStore]}</S.DataStoreName>
                {isConfigured && (
                  <Popover content="This data source is currently configured" placement="right">
                    <S.InfoIcon />
                  </Popover>
                )}
              </S.DataStoreItemContainer>
            }
          />
        );
      })}
    </S.DataStoreListContainer>
  );
};

export default DataStoreSelection;
