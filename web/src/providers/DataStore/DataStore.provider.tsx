import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';

import {
  useTestConnectionMutation,
  useCreateDataStoreMutation,
  useUpdateDataStoreMutation,
  useDeleteDataStoreMutation,
} from 'redux/apis/TraceTest.api';
import {SupportedDataStores, TConnectionResult, TDataStore, TDraftDataStore} from 'types/Config.types';
import DataStoreService from 'services/DataStore.service';
import ConnectionResult from 'models/ConnectionResult.model';
import useTestConnectionNotification from './hooks/useTestConnectionNotification';
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
  const {showNotification, contextHolder} = useTestConnectionNotification();
  const {onOpen} = useConfirmationModal();

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

      if (draft.dataStoreType === SupportedDataStores.OtelCollector) {
        return showNotification(ConnectionResult({}), draft.dataStoreType);
      }

      try {
        const result = await testConnection(dataStore!).unwrap();
        showNotification(result, draft.dataStoreType!);
      } catch (err) {
        showNotification(err as TConnectionResult, draft.dataStoreType!);
      }
    },
    [showNotification, testConnection]
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
