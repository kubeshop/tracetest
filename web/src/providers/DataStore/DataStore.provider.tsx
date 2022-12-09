import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {useTestConnectionMutation, useUpdateDatastoreConfigMutation} from 'redux/apis/TraceTest.api';
import {TDraftDataStore} from 'types/Config.types';
import DataStoreService from 'services/DataStore.service';
import {notification} from 'antd';
import {useTheme} from 'styled-components';

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
  const [api, contextHolder] = notification.useNotification();
  const {
    notification: {success, error},
  } = useTheme();

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
      const {successful, errorMessage = ''} = await testConnection(dataStore!).unwrap();

      if (successful)
        return api.success({
          message: 'Connection is setup',
          description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor',
          ...success,
        });

      api.error({
        message: 'Connection is not setup',
        description: errorMessage,
        ...error,
      });
    },
    [api, error, success, testConnection]
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
