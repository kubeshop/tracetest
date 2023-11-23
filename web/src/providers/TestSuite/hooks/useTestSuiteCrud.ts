import {useCallback} from 'react';
import {TVariableSetValue} from 'models/VariableSet.model';
import RunError from 'models/RunError.model';
import TestSuite from 'models/TestSuite.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useVariableSet} from 'providers/VariableSet/VariableSet.provider';
import {useMissingVariablesModal} from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import TracetestAPI from 'redux/apis/Tracetest';
import {RunErrorTypes} from 'types/TestRun.types';
import {TDraftTestSuite} from 'types/TestSuite.types';
import TestSuiteService from 'services/TestSuite.service';

const {
  useEditTestSuiteMutation,
  useRunTestSuiteMutation,
  useLazyGetTestSuiteVersionByIdQuery,
  useDeleteTestSuiteByIdMutation,
  useCreateTestSuiteMutation,
} = TracetestAPI.instance;

const useTestSuiteCrud = () => {
  const {navigate} = useDashboard();
  const [editTestSuite, {isLoading: isTestSuiteEditLoading}] = useEditTestSuiteMutation();
  const [runTestSuiteAction, {isLoading: isLoadingRunTestSuite}] = useRunTestSuiteMutation();
  const [getTestSuite] = useLazyGetTestSuiteVersionByIdQuery();
  const [deleteTestSuiteAction] = useDeleteTestSuiteByIdMutation();
  const [createAction, {isLoading: isLoadingCreate}] = useCreateTestSuiteMutation();
  const isEditLoading = isTestSuiteEditLoading || isLoadingRunTestSuite;
  const {selectedVariableSet} = useVariableSet();
  const {onOpen} = useMissingVariablesModal();

  const runTestSuite = useCallback(
    async (suite: TestSuite, runId?: number, variableSetId = selectedVariableSet?.id) => {
      const {fullSteps: testList} = await getTestSuite({
        testSuiteId: suite.id,
        version: suite.version,
      }).unwrap();

      const run = async (variables: TVariableSetValue[] = []) => {
        try {
          const {id} = await runTestSuiteAction({testSuiteId: suite.id, variableSetId, variables}).unwrap();

          navigate(`/testsuite/${suite.id}/run/${id}`);
        } catch (error) {
          const {type, missingVariables} = error as RunError;
          if (type === RunErrorTypes.MissingVariables)
            onOpen({
              name: suite.name,
              missingVariables,
              testList,
              onSubmit(missing) {
                run(missing);
              },
            });
          else throw error;
        }
      };

      run();
    },
    [getTestSuite, navigate, onOpen, runTestSuiteAction, selectedVariableSet?.id]
  );

  const edit = useCallback(
    async (suite: TestSuite, draft: TDraftTestSuite) => {
      await editTestSuite({testSuiteId: suite.id, draft}).unwrap();

      runTestSuite(suite);
    },
    [editTestSuite, runTestSuite]
  );

  const deleteTestSuite = useCallback(
    (testSuiteId: string) => {
      deleteTestSuiteAction({testSuiteId});

      navigate('/');
    },
    [deleteTestSuiteAction, navigate]
  );

  const create = useCallback(
    async (values: TDraftTestSuite) => {
      const suite = await createAction(values).unwrap();
      runTestSuite(suite);
    },
    [createAction, runTestSuite]
  );

  const duplicate = useCallback(
    async (suite: TestSuite) => {
      const draft = TestSuiteService.getDuplicatedDraftTestSuite(suite, `${suite.name} (copy)`);

      await create(draft);
    },
    [create]
  );

  return {
    edit,
    runTestSuite,
    deleteTestSuite,
    duplicate,
    create,
    isEditLoading,
    isLoadingCreate,
  };
};

export default useTestSuiteCrud;
