import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {useUpdateDatastoreConfigMutation} from 'redux/apis/TraceTest.api';
import {TDraftDataStore} from 'types/Config.types';
import SetupConfigService from 'services/DataStore.service';

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

  const onSaveConfig = useCallback(
    async (draft: TDraftDataStore) => {
      const configRequest = await SetupConfigService.getRequest(draft);
      const config = await updateConfig(configRequest).unwrap();

      console.log('@@config', config);
    },
    [updateConfig]
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
