import {createContext, useCallback, useContext, useMemo} from 'react';
import {noop} from 'lodash';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {IPlugin} from 'types/Plugins.types';
import {initialState, setDraftTest, setPlugin, reset, setIsFormValid} from 'redux/slices/CreateTest.slice';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import CreateTestSelectors from 'selectors/CreateTest.selectors';
import TracetestAPI from 'redux/apis/Tracetest';
import {ICreateTestState, TDraftTest} from 'types/Test.types';
import TestService from 'services/Test.service';
import {Plugins} from 'constants/Plugins.constants';
import useTestCrud from '../Test/hooks/useTestCrud';

const {useCreateTestMutation} = TracetestAPI.instance;

interface IContext extends ICreateTestState {
  isLoading: boolean;
  plugin: IPlugin;
  onCreateTest(draftTest: TDraftTest, plugin: IPlugin): void;
  onUpdateDraftTest(draftTest: TDraftTest): void;
  onUpdatePlugin(plugin: IPlugin): void;
  onReset(): void;
  onIsFormValid(isValid: boolean): void;
}

export const Context = createContext<IContext>({
  ...initialState,
  isLoading: false,
  plugin: Plugins.REST,
  onCreateTest: noop,
  onUpdateDraftTest: noop,
  onUpdatePlugin: noop,
  onReset: noop,
  onIsFormValid: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useCreateTest = () => useContext(Context);

const CreateTestProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const [createTest, {isLoading: isLoadingCreateTest}] = useCreateTestMutation();
  const {runTest, isEditLoading} = useTestCrud();
  const {demos} = useSettingsValues();

  const draftTest = useAppSelector(CreateTestSelectors.selectDraftTest);
  const plugin = useAppSelector(state => CreateTestSelectors.selectPlugin(state, demos));
  const isFormValid = useAppSelector(CreateTestSelectors.selectIsFormValid);

  const onCreateTest = useCallback(
    async (draft: TDraftTest, p: IPlugin) => {
      const rawTest = await TestService.getRequest(p, draft);
      const test = await createTest(rawTest).unwrap();
      runTest({test});
    },
    [createTest, runTest]
  );

  const onUpdateDraftTest = useCallback(
    (update: TDraftTest) => {
      dispatch(setDraftTest({draftTest: update}));
    },
    [dispatch]
  );

  const onUpdatePlugin = useCallback(
    (newPlugin: IPlugin) => {
      dispatch(setPlugin({plugin: newPlugin}));
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
      draftTest,
      pluginName: plugin.name,
      plugin,
      isLoading: isLoadingCreateTest || isEditLoading,
      isFormValid,
      onCreateTest,
      onUpdateDraftTest,
      onUpdatePlugin,
      onReset,
      onIsFormValid,
    }),
    [
      draftTest,
      plugin,
      isLoadingCreateTest,
      isEditLoading,
      isFormValid,
      onCreateTest,
      onUpdateDraftTest,
      onUpdatePlugin,
      onReset,
      onIsFormValid,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CreateTestProvider;
