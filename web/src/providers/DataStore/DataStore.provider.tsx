import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {notification} from 'antd';
import {useTheme} from 'styled-components';
import {useTestConnectionMutation, useUpdateDatastoreConfigMutation} from 'redux/apis/TraceTest.api';
import {TDraftDataStore} from 'types/Config.types';
import DataStoreService from 'services/DataStore.service';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';

interface IContext {
  isFormValid: boolean;
  isLoading: boolean;
  isTestConnectionLoading: boolean;
  onDeleteConfig(): void;
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
  onDeleteConfig: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useSetupConfig = () => useContext(Context);

const DataStoreProvider = ({children}: IProps) => {
  const [updateConfig, {isLoading}] = useUpdateDatastoreConfigMutation();
  const [testConnection, {isLoading: isTestConnectionLoading}] = useTestConnectionMutation();
  const [isFormValid, setIsFormValid] = useState(false);
  const {onOpen} = useConfirmationModal();
  const [api, contextHolder] = notification.useNotification();
  const {
    notification: {success, error},
  } = useTheme();

  const onSaveConfig = useCallback(
    async (draft: TDraftDataStore) => {
      onOpen({
        title: 'Tracetest needs to do a quick restart to use this new configuration.',
        heading: 'Save Confirmation',
        okText: 'Save & Restart',
        onConfirm: async () => {
          const update = await DataStoreService.getRequest(draft);
          console.log('@@saving draft', draft, update);
          // const config = await updateConfig(update).unwrap();
        },
      });
    },
    [onOpen]
  );

  const onDeleteConfig = useCallback(async () => {
    onOpen({
      title: 'Tracetest needs to do a quick restart to use this new configuration.',
      heading: 'Save Confirmation',
      okText: 'Save & Restart',
      onConfirm: async () => {
        const deleteRequest = await DataStoreService.getDeleteRequest();
        console.log('@@deleting', deleteRequest);
        // const config = await updateConfig(configRequest).unwrap();
      },
    });
  }, [onOpen]);

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
      onDeleteConfig,
    }),
    [isLoading, isFormValid, isTestConnectionLoading, onSaveConfig, onIsFormValid, onTestConnection, onDeleteConfig]
  );

  return (
    <>
      {contextHolder}
      <Context.Provider value={value}>{children}</Context.Provider>
    </>
  );
};

export default DataStoreProvider;
