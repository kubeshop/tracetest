import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';

import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import TracetestAPI from 'redux/apis/Tracetest';
import DataStoreService from 'services/DataStore.service';
import {useContactUsModal} from 'components/ContactUs';
import {SupportedDataStores, TConnectionResult, TDraftDataStore} from 'types/DataStore.types';
import DataStore from 'models/DataStore.model';
import OTLPTestConnectionResponse from 'models/OTLPTestConnectionResponse.model';
import useDataStoreNotification from './hooks/useDataStoreNotification';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';
import {useSettingsValues} from '../SettingsValues/SettingsValues.provider';
import {useWizard} from '../Wizard';

interface IContext {
  isFormValid: boolean;
  isLoading: boolean;
  isTestConnectionLoading: boolean;
  isTestConnectionSuccessful: boolean;
  testConnectionResponse?: TConnectionResult;
  otlpTestConnectionResponse?: OTLPTestConnectionResponse;
  onSaveConfig(draft: TDraftDataStore, defaultDataStore: DataStore, onAfterSave?: () => void): void;
  onTestConnection(draft: TDraftDataStore, defaultDataStore: DataStore): void;
  onSetOtlpTestConnectionResponse(response?: OTLPTestConnectionResponse): void;
  onIsFormValid(isValid: boolean): void;
  resetTestConnection(): void;
}

export const Context = createContext<IContext>({
  isFormValid: false,
  isLoading: false,
  isTestConnectionLoading: false,
  isTestConnectionSuccessful: false,
  onSaveConfig: noop,
  onIsFormValid: noop,
  onTestConnection: noop,
  resetTestConnection: noop,
  onSetOtlpTestConnectionResponse: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useDataStore = () => useContext(Context);

const DataStoreProvider = ({children}: IProps) => {
  const {useTestConnectionMutation, useUpdateDataStoreMutation} = TracetestAPI.instance;
  const {isFetching} = useSettingsValues();
  const {onCompleteStep} = useWizard();
  const [updateDataStore, {isLoading: isLoadingUpdate}] = useUpdateDataStoreMutation();
  const [
    testConnection,
    {isLoading: isTestConnectionLoading, data: testConnectionResponse, reset: resetTestConnection},
  ] = useTestConnectionMutation();
  const [otlpTestConnectionResponse, setOtlpTestConnectionResponse] = useState<
    OTLPTestConnectionResponse | undefined
  >();

  const [isFormValid, setIsFormValid] = useState(false);
  const {showSuccessNotification, showTestConnectionNotification} = useDataStoreNotification();
  const {onOpen} = useConfirmationModal();
  const [connectionTries, setConnectionTries] = useState(0);
  const {onOpen: onContactUsOpen} = useContactUsModal();

  const onSaveConfig = useCallback(
    async (draft: TDraftDataStore, defaultDataStore: DataStore, onAfterSave: () => void = noop) => {
      const warningMessage =
        !!defaultDataStore.id && draft.dataStoreType !== defaultDataStore.type
          ? `Saving will delete your previous configuration of the ${
              SupportedDataStoresToName[defaultDataStore.type || SupportedDataStores.JAEGER]
            } Tracing Backend.`
          : '';

      onOpen({
        title: (
          <>
            <p>Are you sure you want to save this Tracing Backend configuration?</p>
            <p>{warningMessage}</p>
          </>
        ),
        heading: 'Save Confirmation',
        okText: 'Save',
        onConfirm: async () => {
          const dataStore = await DataStoreService.getRequest(draft, defaultDataStore);
          await updateDataStore({dataStore}).unwrap();
          await onCompleteStep('tracing_backend');
          showSuccessNotification();
          onAfterSave();
        },
      });
    },
    [onCompleteStep, onOpen, showSuccessNotification, updateDataStore]
  );

  const onResetTestConnection = useCallback(() => {
    resetTestConnection();
  }, [resetTestConnection]);

  const onIsFormValid = useCallback((isValid: boolean) => {
    setIsFormValid(isValid);
  }, []);

  const onTestConnection = useCallback(
    async (draft: TDraftDataStore, defaultDataStore: DataStore) => {
      const dataStore = await DataStoreService.getRequest(draft, defaultDataStore);

      if (DataStoreService.getIsOtlpBased(draft)) {
        return;
      }

      try {
        const result = await testConnection(dataStore.spec!).unwrap();
        showTestConnectionNotification(result, draft.dataStoreType!);
        setConnectionTries(0);
      } catch (err) {
        setConnectionTries(prev => prev + 1);
        showTestConnectionNotification(err as TConnectionResult, draft.dataStoreType!);
        if (connectionTries + 1 === 3) {
          onContactUsOpen();
        }
      }
    },
    [connectionTries, onContactUsOpen, showTestConnectionNotification, testConnection]
  );

  const value = useMemo<IContext>(
    () => ({
      isLoading: isLoadingUpdate,
      isFormValid,
      isTestConnectionLoading,
      isTestConnectionSuccessful: testConnectionResponse?.allPassed || !!otlpTestConnectionResponse?.spanCount,
      onSaveConfig,
      onIsFormValid,
      onTestConnection,
      testConnectionResponse,
      otlpTestConnectionResponse,
      resetTestConnection: onResetTestConnection,
      onSetOtlpTestConnectionResponse: setOtlpTestConnectionResponse,
    }),
    [
      isFormValid,
      isLoadingUpdate,
      isTestConnectionLoading,
      onIsFormValid,
      onResetTestConnection,
      onSaveConfig,
      onTestConnection,
      otlpTestConnectionResponse,
      testConnectionResponse,
    ]
  );

  return <Context.Provider value={value}>{isFetching ? null : children}</Context.Provider>;
};

export default DataStoreProvider;
