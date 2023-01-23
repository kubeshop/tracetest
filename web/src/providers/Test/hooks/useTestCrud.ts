import {useCallback} from 'react';
import {useMatch, useNavigate} from 'react-router-dom';
import {useAppDispatch} from 'redux/hooks';
import {reset} from 'redux/slices/TestSpecs.slice';
import {TDraftTest, TTest} from 'types/Test.types';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import TestService from 'services/Test.service';
import {TEnvironmentValue} from 'types/Environment.types';
import {useEditTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {useMissingVariablesModal} from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import {RunErrorTypes, TRunError} from 'types/TestRun.types';

const useTestCrud = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const {updateIsInitialized} = useTestSpecs();
  const [editTest, {isLoading: isLoadingEditTest}] = useEditTestMutation();
  const [runTestAction, {isLoading: isLoadingRunTest}] = useRunTestMutation();
  const isEditLoading = isLoadingEditTest || isLoadingRunTest;
  const match = useMatch('/test/:testId/run/:runId/:mode');
  const {selectedEnvironment} = useEnvironment();
  const {onOpen} = useMissingVariablesModal();

  const runTest = useCallback(
    async (test: TTest, runId?: string, environmentId = selectedEnvironment?.id) => {
      const run = async (variables: TEnvironmentValue[] = []) => {
        try {
          TestAnalyticsService.onRunTest();
          const {id} = await runTestAction({testId: test.id, environmentId, variables}).unwrap();
          dispatch(reset());

          const mode = match?.params.mode || 'trigger';
          navigate(`/test/${test.id}/run/${id}/${mode}`);
        } catch (error) {
          const {type, missingVariables} = error as TRunError;
          if (type === RunErrorTypes.MissingVariables)
            onOpen({
              name: test.name,
              missingVariables,
              onSubmit(missing) {
                run(missing);
              },
            });
          else throw error;
        }
      };

      run();
    },
    [dispatch, match?.params.mode, navigate, onOpen, runTestAction, selectedEnvironment?.id]
  );

  const edit = useCallback(
    async (test: TTest, draft: TDraftTest) => {
      const {id: testId, trigger} = test;
      updateIsInitialized(false);
      const plugin = TriggerTypeToPlugin[trigger.type || TriggerTypes.http];
      const rawTest = await TestService.getRequest(plugin, draft, test);

      await editTest({
        test: rawTest,
        testId,
      }).unwrap();

      runTest(test);
    },
    [editTest, runTest, updateIsInitialized]
  );

  return {
    edit,
    runTest,
    isEditLoading,
    isLoadingRunTest,
  };
};

export default useTestCrud;
