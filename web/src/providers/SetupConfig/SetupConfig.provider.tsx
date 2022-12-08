import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';
import {useUpdateConfigMutation} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {initialState, reset, setDraftConfig, setIsFormValid} from 'redux/setupConfig/slice';
import SetupConfigSelectors from 'redux/setupConfig/selectors';
import {ISetupConfigState, TDraftConfig} from 'types/Config.types';
// import SetupConfigService from 'services/SetupConfig.service';

interface IContext extends ISetupConfigState {
  isLoading: boolean;
  onSaveConfig(draft: TDraftConfig): void;
  onUpdateDraft(draft: TDraftConfig): void;
  onReset(): void;
  onIsFormValid(isValid: boolean): void;
}

export const Context = createContext<IContext>({
  ...initialState,
  isLoading: false,
  onSaveConfig: noop,
  onUpdateDraft: noop,
  onIsFormValid: noop,
  onReset: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useSetupConfig = () => useContext(Context);

const SetupConfigProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const [updateConfig, {isLoading}] = useUpdateConfigMutation();

  const draftConfig = useAppSelector(SetupConfigSelectors.selectDraftConfig);
  const isFormValid = useAppSelector(SetupConfigSelectors.selectIsFormValid);

  const onSaveConfig = useCallback(
    async (draft: TDraftConfig) => {
      console.log('@@saving', draft);

      // const configRequest = await SetupConfigService.getRequest(draft);
      // const config = await updateConfig(configRequest).unwrap();
    },
    [updateConfig]
  );

  const onUpdateDraft = useCallback(
    (update: TDraftConfig) => {
      dispatch(setDraftConfig({draftConfig: update}));
    },
    [dispatch]
  );

  const onReset = useCallback(() => {
    dispatch(reset());
  }, [dispatch]);

  const onIsFormValid = useCallback(
    (isValid: boolean) => {
      dispatch(setIsFormValid({isValid}));
    },
    [dispatch]
  );

  const value = useMemo<IContext>(
    () => ({
      draftConfig,
      isLoading,
      isFormValid,
      onSaveConfig,
      onUpdateDraft,
      onReset,
      onIsFormValid,
    }),
    [draftConfig, isLoading, isFormValid, onSaveConfig, onUpdateDraft, onReset, onIsFormValid]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SetupConfigProvider;
