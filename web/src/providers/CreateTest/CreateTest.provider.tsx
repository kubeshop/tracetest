import {noop} from 'lodash';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import TestService from 'services/Test.service';
import {IPlugin} from 'types/Plugins.types';
import {TDraftTest} from 'types/Test.types';

const {useCreateTestMutation} = TracetestAPI.instance;

interface IContext {
  initialValues: TDraftTest;
  isLoading: boolean;
  onCreateTest(draftTest: TDraftTest, plugin: IPlugin): void;
  onInitialValues(draftTest: TDraftTest): void;
}

export const Context = createContext<IContext>({
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
      initialValues,
      isLoading: isLoadingCreateTest || isLoadingEditTest,
      onCreateTest,
      onInitialValues,
    }),
    [initialValues, isLoadingCreateTest, isLoadingEditTest, onCreateTest, onInitialValues]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CreateTestProvider;
