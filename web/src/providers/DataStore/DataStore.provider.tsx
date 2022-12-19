import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {useTestConnectionMutation, useUpdateDatastoreConfigMutation} from 'redux/apis/TraceTest.api';
import {TConnectionResult, TDraftDataStore} from 'types/Config.types';
import DataStoreService from 'services/DataStore.service';
import useTestConnectionNotification from './hooks/useTestConnectionNotification';

interface IContext {
  isFormValid: boolean;
  isLoading: boolean;
  isTestConnectionLoading: boolean;
  onSaveConfig(draft: TDraftDataStore): void;
  onTestConnection(draft: TDraftDataStore): void;
  onIsFormValid(isValid: boolean): void;
}

export const Context = createContext<IContext>({
  isFormValid: false,
  isLoading: false,
  isTestConnectionLoading: false,
  onSaveConfig: noop,
  onIsFormValid: noop,
  onTestConnection: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useSetupConfig = () => useContext(Context);

const DataStoreProvider = ({children}: IProps) => {
  const [updateConfig, {isLoading}] = useUpdateDatastoreConfigMutation();
  const [testConnection, {isLoading: isTestConnectionLoading}] = useTestConnectionMutation();
  const [isFormValid, setIsFormValid] = useState(false);
  const {showNotification, contextHolder} = useTestConnectionNotification();

  const onSaveConfig = useCallback(async (draft: TDraftDataStore) => {
    const configRequest = await DataStoreService.getRequest(draft);
    console.log('@@saving draft', draft, configRequest);
    // const config = await updateConfig(configRequest).unwrap();
  }, []);

  const onIsFormValid = useCallback((isValid: boolean) => {
    setIsFormValid(isValid);
  }, []);

  const onTestConnection = useCallback(
    async (draft: TDraftDataStore) => {
      const {dataStores: [dataStore] = []} = await DataStoreService.getRequest(draft);

      try {
        const result = await testConnection(dataStore!).unwrap();
        showNotification(result);
      } catch (err) {
        showNotification(err as TConnectionResult);
      }
    },
    [showNotification, testConnection]
  );

  const value = useMemo<IContext>(
    () => ({
      isLoading,
      isFormValid,
      isTestConnectionLoading,
      onSaveConfig,
      onIsFormValid,
      onTestConnection,
    }),
    [isLoading, isFormValid, isTestConnectionLoading, onSaveConfig, onIsFormValid, onTestConnection]
  );

  return (
    <>
      {contextHolder}
      <Context.Provider value={value}>{children}</Context.Provider>
    </>
  );
};

export default DataStoreProvider;
