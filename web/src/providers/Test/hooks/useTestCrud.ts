import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useAppDispatch} from 'redux/hooks';
import {reset} from 'redux/slices/TestDefinition.slice';
import {TDraftTest, TTest} from 'types/Test.types';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import TestService from 'services/Test.service';
import {useEditTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';

const useTestCrud = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const {updateIsInitialized} = useTestDefinition();
  const [editTest, {isLoading: isLoadingEditTest}] = useEditTestMutation();
  const [runTestAction, {isLoading: isLoadingRunTest}] = useRunTestMutation();
  const isEditLoading = isLoadingEditTest || isLoadingRunTest;

  const runTest = useCallback(
    async (testId: string) => {
      TestAnalyticsService.onRunTest();
      const run = await runTestAction({testId}).unwrap();
      dispatch(reset());

      navigate(`/test/${testId}/run/${run.id}`);
    },
    [dispatch, navigate, runTestAction]
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

      runTest(testId);
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
