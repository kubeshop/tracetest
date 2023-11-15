import {getDemoByPluginMap} from 'constants/Demo.constants';
import {noop} from 'lodash';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import TestService from 'services/Test.service';
import {IPlugin} from 'types/Plugins.types';
import {TDraftTest} from 'types/Test.types';

const {useCreateTestMutation} = TracetestAPI.instance;

interface IContext {
  demoList: TDraftTest[];
  initialValues: TDraftTest;
  isLoading: boolean;
  onCreateTest(draftTest: TDraftTest, plugin: IPlugin): void;
  onInitialValues(draftTest: TDraftTest): void;
}

export const Context = createContext<IContext>({
  demoList: [],
  initialValues: {},
  isLoading: false,
  onCreateTest: noop,
  onInitialValues: noop,
});

export const useCreateTest = () => useContext(Context);

interface IProps {
  children: React.ReactNode;
}

const CreateTestProvider = ({children}: IProps) => {
  const [initialValues, setInitialValues] = useState<TDraftTest>({name: 'Untitled'});
  const [createTest, {isLoading: isLoadingCreateTest}] = useCreateTestMutation();
  const {runTest, isEditLoading: isLoadingEditTest} = useTestCrud();

  // TODO: this is a hack to get the demo list for REST plugin
  const {demos} = useSettingsValues();
  const demoByPluginMap = getDemoByPluginMap(demos);
  const demoList = demoByPluginMap['REST'];

  const onCreateTest = useCallback(
    async (draft: TDraftTest, plugin: IPlugin) => {
      const rawTest = await TestService.getRequest(plugin, draft);
      const test = await createTest(rawTest).unwrap();
      runTest({test});
    },
    [createTest, runTest]
  );

  const onInitialValues = useCallback(values => setInitialValues(values), []);

  const value = useMemo<IContext>(
    () => ({
      demoList,
      initialValues,
      isLoading: isLoadingCreateTest || isLoadingEditTest,
      onCreateTest,
      onInitialValues,
    }),
    [demoList, initialValues, isLoadingCreateTest, isLoadingEditTest, onCreateTest, onInitialValues]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CreateTestProvider;
