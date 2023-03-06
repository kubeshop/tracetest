import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';

import {NoTestConnectionDataStoreList, SupportedDataStoresToName} from 'constants/DataStore.constants';
import ConnectionResult from 'models/ConnectionResult.model';
import {
  useTestConnectionMutation,
  useCreateDataStoreMutation,
  useUpdateDataStoreMutation,
  useDeleteDataStoreMutation,
} from 'redux/apis/TraceTest.api';
import DataStoreService from 'services/DataStore.service';
import {SupportedDataStores, TConnectionResult, TDraftDataStore} from 'types/DataStore.types';
import DataStore from 'models/DataStore.model';
import useDataStoreNotification from './hooks/useDataStoreNotification';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';
import {useDataStoreConfig} from '../DataStoreConfig/DataStoreConfig.provider';

interface IContext {
  isFormValid: boolean;
  isLoading: boolean;
  isTestConnectionLoading: boolean;
  onDeleteConfig(defaultDataStore: DataStore): void;
  onSaveConfig(draft: TDraftDataStore, defaultDataStore: DataStore): void;
  onTestConnection(draft: TDraftDataStore, defaultDataStore: DataStore): void;
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

export const useDataStore = () => useContext(Context);

const DataStoreProvider = ({children}: IProps) => {
  const {isFetching} = useDataStoreConfig();
  const [createDataStore, {isLoading: isLoadingCreate}] = useCreateDataStoreMutation();
  const [updateDataStore, {isLoading: isLoadingUpdate}] = useUpdateDataStoreMutation();
  const [deleteDataStore] = useDeleteDataStoreMutation();
  const [testConnection, {isLoading: isTestConnectionLoading}] = useTestConnectionMutation();
  const [isFormValid, setIsFormValid] = useState(false);
  const {showSuccessNotification, showTestConnectionNotification} = useDataStoreNotification();
  const {onOpen} = useConfirmationModal();

  const onSaveConfig = useCallback(
    async (draft: TDraftDataStore, defaultDataStore: DataStore) => {
      const warningMessage =
        !!defaultDataStore.id && draft.dataStoreType !== defaultDataStore.type
          ? `Saving will delete your previous configuration of the ${
              SupportedDataStoresToName[defaultDataStore.type || SupportedDataStores.JAEGER]
            } data store`
          : '';

      onOpen({
        title: (
          <>
            <p>Are you sure you want to save this Data Store configuration?</p>
            <p>{warningMessage}</p>
          </>
        ),
        heading: 'Save Confirmation',
        okText: 'Save',
        onConfirm: async () => {
          const dataStore = await DataStoreService.getRequest(draft, defaultDataStore);
          if (dataStore.id) {
            await updateDataStore({dataStore, dataStoreId: dataStore.id}).unwrap();
          } else {
            await createDataStore(dataStore).unwrap();
          }
          showSuccessNotification();
        },
      });
    },
    [createDataStore, onOpen, showSuccessNotification, updateDataStore]
  );

  const onDeleteConfig = useCallback(
    async (defaultDataStore: DataStore) => {
      onOpen({
        title:
          "Tracetest will remove the trace data store configuration information and enter the 'No-Tracing Mode'. You can still run tests against the responses until you configure a new trace data store.",
        heading: 'Save Confirmation',
        okText: 'Save',
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
    async (draft: TDraftDataStore, defaultDataStore: DataStore) => {
      const dataStore = await DataStoreService.getRequest(draft, defaultDataStore);

      if (NoTestConnectionDataStoreList.includes(draft.dataStoreType!)) {
        return showTestConnectionNotification(ConnectionResult({}), draft.dataStoreType!);
      }

      try {
        const result = await testConnection(dataStore!).unwrap();
        showTestConnectionNotification(result, draft.dataStoreType!);
      } catch (err) {
        showTestConnectionNotification(err as TConnectionResult, draft.dataStoreType!);
      }
    },
    [showTestConnectionNotification, testConnection]
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

  return <Context.Provider value={value}>{isFetching ? null : children}</Context.Provider>;
};

export default DataStoreProvider;
