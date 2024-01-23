import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import * as S from './TracingBackend.styled';

const Tab = () => {
  const {dataStoreConfig} = useSettingsValues();

  return (
    <div>
      <S.TabText>{SupportedDataStoresToName[dataStoreConfig.defaultDataStore.type]} connected</S.TabText>
    </div>
  );
};

export default Tab;
