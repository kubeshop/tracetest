import {Popover, Tabs} from 'antd';
import {useCallback} from 'react';
import {noop} from 'lodash';
import {useTheme} from 'styled-components';
import {ConfigMode, SupportedDataStores} from 'types/DataStore.types';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import {Flag, useCustomization} from 'providers/Customization';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';

import DataStoreIcon from '../../DataStoreIcon/DataStoreIcon';
import * as S from './DataStoreSelection.styled';

interface IProps {
  value?: SupportedDataStores;
  onChange?(value: SupportedDataStores): void;
}

const supportedDataStoreList = Object.values(SupportedDataStores);

const DataStoreSelection = ({onChange = noop, value = SupportedDataStores.JAEGER}: IProps) => {
  const {getFlag} = useCustomization();
  const isLocalModeEnabled = getFlag(Flag.IsLocalModeEnabled);
  const {
    color: {text, primary},
  } = useTheme();
  const {dataStoreConfig} = useSettingsValues();
  const configuredDataStoreType = dataStoreConfig.defaultDataStore.type;

  const handleChange = useCallback(
    dataStore => {
      const isDisabled = isLocalModeEnabled && dataStore !== SupportedDataStores.Agent;

      if (!isDisabled) onChange(dataStore);
    },
    [isLocalModeEnabled, onChange]
  );

  return (
    <S.DataStoreListContainer tabPosition="left" onChange={handleChange}>
      {supportedDataStoreList.map(dataStore => {
        if (dataStore === SupportedDataStores.Agent && !getFlag(Flag.IsAgentDataStoreEnabled)) {
          return null;
        }

        const isSelected = value === dataStore;
        const isConfigured = configuredDataStoreType === dataStore && dataStoreConfig.mode === ConfigMode.READY;
        const isDisabled = isLocalModeEnabled && dataStore !== SupportedDataStores.Agent;

        return (
          <Tabs.TabPane
            key={dataStore}
            tab={
              <S.DataStoreItemContainer $isDisabled={isDisabled} $isSelected={isSelected} key={dataStore}>
                <DataStoreIcon dataStoreType={dataStore} color={isSelected ? primary : text} width="22" height="22" />

                {isDisabled ? (
                  <Popover
                    content={
                      <div>
                        In localMode only the Agent Tracing Backend can be used. <br /> If you want to connect to a
                        different Tracing Backend <br /> please create a new environment
                      </div>
                    }
                    placement="right"
                  >
                    <S.DataStoreName $isSelected={isSelected}>{SupportedDataStoresToName[dataStore]}</S.DataStoreName>
                  </Popover>
                ) : (
                  <S.DataStoreName $isSelected={isSelected}>{SupportedDataStoresToName[dataStore]}</S.DataStoreName>
                )}

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
