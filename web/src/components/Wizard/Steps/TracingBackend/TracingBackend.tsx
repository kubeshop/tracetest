import {useCallback, useEffect, useState} from 'react';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import DataStore, {fromType} from 'models/DataStore.model';
import WizardAnalytics from 'services/Analytics/Wizard.service';
import Selector from './Selector';
import Configuration from './Configuration';

const TracingBackend = () => {
  const {
    dataStoreConfig: {defaultDataStore},
  } = useSettingsValues();
  const {resetTestConnection} = useDataStore();
  const [selectedDataStore, setSelectedDataStore] = useState<DataStore | undefined>(defaultDataStore);

  const handleOnSelect = useCallback(
    type => {
      setSelectedDataStore(type === defaultDataStore.type ? defaultDataStore : fromType(type));
      WizardAnalytics.onTracingBackendTypeSelect(type);
    },
    [defaultDataStore]
  );

  useEffect(() => {
    resetTestConnection();
  }, [selectedDataStore?.type]);

  if (!selectedDataStore) return <Selector onSelect={handleOnSelect} />;

  return <Configuration dataStore={selectedDataStore} onBack={() => setSelectedDataStore(undefined)} />;
};

export default TracingBackend;
