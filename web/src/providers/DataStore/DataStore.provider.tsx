import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {notification} from 'antd';
import {useTheme} from 'styled-components';
import {
  useTestConnectionMutation,
  useCreateDataStoreMutation,
  useUpdateDataStoreMutation,
  useDeleteDataStoreMutation,
} from 'redux/apis/TraceTest.api';
import {TDataStore, TDraftDataStore} from 'types/Config.types';
import DataStoreService from 'services/DataStore.service';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';
import {useDataStoreConfig} from '../DataStoreConfig/DataStoreConfig.provider';

interface IContext {
  isFormValid: boolean;
  isLoading: boolean;
  isTestConnectionLoading: boolean;
  onDeleteConfig(defaultDataStore: TDataStore): void;
  onSaveConfig(draft: TDraftDataStore, defaultDataStore: TDataStore): void;
  onTestConnection(draft: TDraftDataStore, defaultDataStore: TDataStore): void;
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
  const {isFetching} = useDataStoreConfig();
  const [createDataStore, {isLoading: isLoadingCreate}] = useCreateDataStoreMutation();
  const [updateDataStore, {isLoading: isLoadingUpdate}] = useUpdateDataStoreMutation();
  const [deleteDataStore] = useDeleteDataStoreMutation();
  const [testConnection, {isLoading: isTestConnectionLoading}] = useTestConnectionMutation();
  const [isFormValid, setIsFormValid] = useState(false);
  const {onOpen} = useConfirmationModal();
  const [api, contextHolder] = notification.useNotification();
  const {
    notification: {success, error},
  } = useTheme();

  const onSaveConfig = useCallback(
    async (draft: TDraftDataStore, defaultDataStore: TDataStore) => {
      onOpen({
        title: 'Tracetest needs to do a quick restart to use this new configuration.',
        heading: 'Save Confirmation',
        okText: 'Save & Restart',
        onConfirm: async () => {
          const dataStore = await DataStoreService.getRequest(draft, defaultDataStore);
          if (dataStore.id) {
            await updateDataStore({dataStore, dataStoreId: dataStore.id}).unwrap();
          } else {
            await createDataStore(dataStore).unwrap();
          }
        },
      });
    },
    [createDataStore, onOpen, updateDataStore]
  );

  const onDeleteConfig = useCallback(
    async (defaultDataStore: TDataStore) => {
      onOpen({
        title: 'Tracetest needs to do a quick restart to use this new configuration.',
        heading: 'Save Confirmation',
        okText: 'Save & Restart',
        onConfirm: async () => {
          await deleteDataStore({dataStoreId: defaultDataStore.id}).unwrap();
        },
      });
    },
    [deleteDataStore, onOpen]
  );

  const onIsFormValid = useCallback((isValid: boolean) => {
    setIsFormValid(isValid);
  }, []);

  const onTestConnection = useCallback(
    async (draft: TDraftDataStore, defaultDataStore: TDataStore) => {
      const dataStore = await DataStoreService.getRequest(draft, defaultDataStore);
      const {authentication = {}, connectivity = {}, fetchTraces = {}} = await testConnection(dataStore!).unwrap();

      if (authentication.passed && connectivity.passed && fetchTraces.passed) {
        return api.success({
          message: 'Connection is setup',
          description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor',
          ...success,
        });
      }

      api.error({
        message: 'Connection is not setup',
        description: authentication.error || connectivity.error || fetchTraces.error,
        ...error,
      });
    },
    [api, error, success, testConnection]
  );

  const value = useMemo<IContext>(
    () => ({
      isLoading: isLoadingCreate || isLoadingUpdate,
      isFormValid,
      isTestConnectionLoading,
      onSaveConfig,
      onIsFormValid,
      onTestConnection,
      onDeleteConfig,
    }),
    [
      isLoadingCreate,
      isLoadingUpdate,
      isFormValid,
      isTestConnectionLoading,
      onSaveConfig,
      onIsFormValid,
      onTestConnection,
      onDeleteConfig,
    ]
  );

  return (
    <>
      {contextHolder}
      <Context.Provider value={value}>{isFetching ? null : children}</Context.Provider>
    </>
  );
};

export default DataStoreProvider;
