import {noop} from 'lodash';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {useGetTestByIdQuery, useGetTestVersionByIdQuery} from 'redux/apis/TraceTest.api';
import {TDraftTest, TTest} from 'types/Test.types';
import VersionMismatchModal from 'components/VersionMismatchModal';
import useTestCrud from './hooks/useTestCrud';

interface IContext {
  onEdit(values: TDraftTest): void;
  onRun(): void;
  isLoading: boolean;
  isError: boolean;
  test: TTest;
  latestTest: TTest;
  isLatestVersion: boolean;
  isEditLoading: boolean;
}

export const Context = createContext<IContext>({
  onEdit: noop,
  onRun: noop,
  test: {} as TTest,
  latestTest: {} as TTest,
  isLoading: false,
  isError: false,
  isLatestVersion: true,
  isEditLoading: false,
});

interface IProps {
  testId: string;
  children: React.ReactNode;
}

export const useTest = () => useContext(Context);

const TestProvider = ({children, testId}: IProps) => {
  const [isVersionModalOpen, setIsVersionModalOpen] = useState(false);
  const [draft, setDraft] = useState<TDraftTest>({});
  const [state, setState] = useState<'edit' | 'run'>();
  const {
    run: {testVersion},
  } = useTestRun();
  const {runTest, edit, isEditLoading} = useTestCrud();
  const {
    data: test,
    isLoading: isCurrentLoading,
    isError: isCurrentError,
  } = useGetTestVersionByIdQuery({testId, version: testVersion}, {skip: !testVersion});
  const {data: latestTest, isLoading: isLatestLoading, isError: isLatestError} = useGetTestByIdQuery({testId});

  const isLatestVersion = useMemo(
    () => Boolean(testVersion) && testVersion === latestTest?.version,
    [latestTest?.version, testVersion]
  );
  const isLoading = isLatestLoading || isCurrentLoading;
  const isError = isLatestError || isCurrentError;
  const currentTest = (test || latestTest)!;

  const onEdit = useCallback(
    (values: TDraftTest) => {
      if (isLatestVersion) edit(test!, values);
      else {
        setState('edit');
        setDraft(values);
        setIsVersionModalOpen(true);
      }
    },
    [edit, isLatestVersion, test]
  );

  const onRun = useCallback(() => {
    if (isLatestVersion) runTest(testId);
    else {
      setState('run');
      setIsVersionModalOpen(true);
    }
  }, [isLatestVersion, runTest, testId]);

  const onConfirm = useCallback(() => {
    if (state === 'edit') edit(test!, draft);
    else runTest(testId);

    setIsVersionModalOpen(false);
  }, [draft, edit, runTest, state, test, testId]);

  const value = useMemo<IContext>(
    () => ({
      onEdit,
      onRun,
      isLoading,
      isError,
      test: currentTest,
      latestTest: latestTest!,
      isLatestVersion,
      isEditLoading,
    }),
    [onEdit, onRun, isLoading, isError, currentTest, latestTest, isLatestVersion, isEditLoading]
  );

  return currentTest && latestTest ? (
    <>
      <Context.Provider value={value}>{children}</Context.Provider>
      <VersionMismatchModal
        description={
          state === 'edit'
            ? 'Changing and saving changes will result in a new version that will become the latest.'
            : 'Running the test will use the latest version of the test.'
        }
        currentVersion={currentTest.version}
        isOpen={isVersionModalOpen}
        latestVersion={latestTest.version}
        okText="Run Test"
        onCancel={() => setIsVersionModalOpen(false)}
        onConfirm={onConfirm}
      />
    </>
  ) : (
    <div data-cy="loading test" />
  );
};

export default TestProvider;
