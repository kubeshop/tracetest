import {useCallback} from 'react';
import {noop} from 'lodash';
import {useMatch, useNavigate} from 'react-router-dom';
import {useAppDispatch} from 'redux/hooks';
import {reset} from 'redux/slices/TestSpecs.slice';
import {TDraftTest} from 'types/Test.types';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import TestService from 'services/Test.service';
import {useEditTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {useMissingVariablesModal} from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import {RunErrorTypes} from 'types/TestRun.types';
import {TEnvironmentValue} from 'models/Environment.model';
import Test from 'models/Test.model';
import RunError from 'models/RunError.model';

export type TTestRunRequest = {
  test: Test;
  environmentId?: string;
  variables?: TEnvironmentValue[];
  onCancel?(): void;
};

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
    async ({test, environmentId = selectedEnvironment?.id, variables = [], onCancel = noop}: TTestRunRequest) => {
      const run = async (updatedVars: TEnvironmentValue[] = variables) => {
        try {
          TestAnalyticsService.onRunTest();
          const {id} = await runTestAction({testId: test.id, environmentId, variables: updatedVars}).unwrap();
          dispatch(reset());

          const mode = match?.params.mode || 'trigger';
          navigate(`/test/${test.id}/run/${id}/${mode}`);
        } catch (error) {
          const {type, missingVariables} = error as RunError;
          if (type === RunErrorTypes.MissingVariables)
            onOpen({
              name: test.name,
              missingVariables,
              testList: [test],
              onSubmit(missing) {
                run(missing);
              },
              onCancel,
            });
          else throw error;
        }
      };

      run();
    },
    [dispatch, match?.params.mode, navigate, onOpen, runTestAction, selectedEnvironment?.id]
  );

  const edit = useCallback(
    async (test: Test, draft: TDraftTest) => {
      const {id: testId, trigger} = test;
      updateIsInitialized(false);
      const plugin = TriggerTypeToPlugin[trigger.type || TriggerTypes.http];
      const rawTest = await TestService.getRequest(plugin, draft, test);

      await editTest({
        test: rawTest,
        testId,
      }).unwrap();

      runTest({
        test,
      });
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
