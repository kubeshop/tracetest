import {useCallback, useState} from 'react';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import DataStore, {fromType} from 'models/DataStore.model';
import WizardAnalytics from 'services/Analytics/Wizard.service';
import {SupportedDataStores} from 'types/DataStore.types';
import {IWizardStepComponentProps} from 'types/Wizard.types';
import Selector from './Selector';
import Configuration from './Configuration';

const TracingBackend = ({onNext}: IWizardStepComponentProps) => {
  const {
    dataStoreConfig: {defaultDataStore},
  } = useSettingsValues();
  const [selectedDataStore, setSelectedDataStore] = useState<DataStore | undefined>();

  const handleOnSelect = useCallback(
    type => {
      setSelectedDataStore(type === defaultDataStore.type ? defaultDataStore : fromType(type));
      WizardAnalytics.onTracingBackendTypeSelect(type);
    },
    [defaultDataStore]
  );

  if (!selectedDataStore)
    return <Selector onSelect={handleOnSelect} selectedBackend={defaultDataStore.type as SupportedDataStores} />;

  return <Configuration dataStore={selectedDataStore} onBack={() => setSelectedDataStore(undefined)} onNext={onNext} />;
};

export default TracingBackend;
