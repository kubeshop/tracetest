import {useCallback, useState} from 'react';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import DataStore, {fromType} from 'models/DataStore.model';
import Selector from './Selector';
import Configuration from './Configuration';

const TracingBackend = () => {
  const {
    dataStoreConfig: {defaultDataStore},
  } = useSettingsValues();
  const [selectedDataStore, setSelectedDataStore] = useState<DataStore | undefined>(defaultDataStore);

  const handleOnSelect = useCallback(
    type => {
      setSelectedDataStore(type === defaultDataStore.type ? defaultDataStore : fromType(type));
    },
    [defaultDataStore]
  );

  if (!selectedDataStore) return <Selector onSelect={handleOnSelect} />;

  return <Configuration dataStore={selectedDataStore} onBack={() => setSelectedDataStore(undefined)} />;
};

export default TracingBackend;
