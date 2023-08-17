import {useCallback} from 'react';
import {noop} from 'lodash';
import {useMatch} from 'react-router-dom';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import {TVariableSetValue} from 'models/VariableSet.model';
import RunError from 'models/RunError.model';
import Test from 'models/Test.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useVariableSet} from 'providers/VariableSet';
import {useMissingVariablesModal} from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import TracetestAPI from 'redux/apis/Tracetest';
import {useAppDispatch} from 'redux/hooks';
import {reset} from 'redux/slices/TestSpecs.slice';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import TestService from 'services/Test.service';
import {TDraftTest} from 'types/Test.types';
import {RunErrorTypes} from 'types/TestRun.types';

const {useEditTestMutation, useRunTestMutation} = TracetestAPI.instance;

export type TTestRunRequest = {
  test: Test;
  variableSetId?: string;
  variables?: TVariableSetValue[];
  onCancel?(): void;
};

const useTestCrud = () => {
  const dispatch = useAppDispatch();
  const {navigate} = useDashboard();
  const {updateIsInitialized} = useTestSpecs();
  const [editTest, {isLoading: isLoadingEditTest}] = useEditTestMutation();
  const [runTestAction, {isLoading: isLoadingRunTest}] = useRunTestMutation();
  const isEditLoading = isLoadingEditTest || isLoadingRunTest;
  const match = useMatch('/test/:testId/run/:runId/:mode');
  const {selectedVariableSet} = useVariableSet();
  const {onOpen} = useMissingVariablesModal();

  const runTest = useCallback(
    async ({test, variableSetId = selectedVariableSet?.id, variables = [], onCancel = noop}: TTestRunRequest) => {
      const run = async (updatedVars: TVariableSetValue[] = variables) => {
        try {
          TestAnalyticsService.onRunTest();
          const {id} = await runTestAction({testId: test.id, variableSetId, variables: updatedVars}).unwrap();
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
    [dispatch, match?.params.mode, navigate, onOpen, runTestAction, selectedVariableSet?.id]
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
