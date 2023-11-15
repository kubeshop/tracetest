import {noop} from 'lodash';
import useTestCrud from 'providers/Test/hooks/useTestCrud';
import {createContext, useCallback, useContext, useMemo} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import TestService from 'services/Test.service';
import {IPlugin} from 'types/Plugins.types';
import {TDraftTest} from 'types/Test.types';

const {useCreateTestMutation} = TracetestAPI.instance;

interface IContext {
  isLoading: boolean;
  onCreateTest(draftTest: TDraftTest, plugin: IPlugin): void;
}

export const Context = createContext<IContext>({
  isLoading: false,
  onCreateTest: noop,
});

export const useCreateTest = () => useContext(Context);

interface IProps {
  children: React.ReactNode;
}

const CreateTestProvider = ({children}: IProps) => {
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

  const value = useMemo<IContext>(
    () => ({
      isLoading: isLoadingCreateTest || isLoadingEditTest,
      onCreateTest,
    }),
    [isLoadingCreateTest, isLoadingEditTest, onCreateTest]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CreateTestProvider;
