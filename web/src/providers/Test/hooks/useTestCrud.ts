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
import {useEditTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';

const useTestCrud = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const {updateIsInitialized} = useTestSpecs();
  const [editTest, {isLoading: isLoadingEditTest}] = useEditTestMutation();
  const [runTestAction, {isLoading: isLoadingRunTest}] = useRunTestMutation();
  const isEditLoading = isLoadingEditTest || isLoadingRunTest;
  const match = useMatch('/test/:testId/run/:runId/:mode');

  const runTest = useCallback(
    async (testId: string) => {
      TestAnalyticsService.onRunTest();
      const run = await runTestAction({testId}).unwrap();
      dispatch(reset());

      const mode = match?.params.mode || 'trigger';

      navigate(`/test/${testId}/run/${run.id}/${mode}`);
    },
    [dispatch, match?.params.mode, navigate, runTestAction]
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
