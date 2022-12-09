import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {useUpdateDatastoreConfigMutation} from 'redux/apis/TraceTest.api';
import {TDraftDataStore} from 'types/Config.types';
import DataStoreService from 'services/DataStore.service';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';

interface IContext {
  isFormValid: boolean;
  isLoading: boolean;
  onSaveConfig(draft: TDraftDataStore): void;
  onIsFormValid(isValid: boolean): void;
}

export const Context = createContext<IContext>({
  isFormValid: false,
  isLoading: false,
  onSaveConfig: noop,
  onIsFormValid: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useSetupConfig = () => useContext(Context);

const SetupConfigProvider = ({children}: IProps) => {
  const [updateConfig, {isLoading}] = useUpdateDatastoreConfigMutation();
  const [isFormValid, setIsFormValid] = useState(false);
  const {onOpen} = useConfirmationModal();

  const onSaveConfig = useCallback(
    async (draft: TDraftDataStore) => {
      onOpen({
        title: 'Tracetest is about to be restarted and there will be some downtime. Please confirm to continue.',
        heading: 'Save Confirmation',
        okText: 'Confirm',
        onConfirm: async () => {
          const configRequest = await DataStoreService.getRequest(draft);
          console.log('@@saving draft', draft, configRequest);
          // const config = await updateConfig(configRequest).unwrap();
        },
      });
    },
    [onOpen]
  );

  const onIsFormValid = useCallback((isValid: boolean) => {
    setIsFormValid(isValid);
  }, []);

  const value = useMemo<IContext>(
    () => ({
      isLoading,
      isFormValid,
      onSaveConfig,
      onIsFormValid,
    }),
    [isLoading, isFormValid, onSaveConfig, onIsFormValid]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SetupConfigProvider;
