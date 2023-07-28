import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {useGetTestByIdQuery, useGetTestVersionByIdQuery} from 'redux/apis/Tracetest';
import {TDraftTest} from 'types/Test.types';
import VersionMismatchModal from 'components/VersionMismatchModal';
import TestService from 'services/Test.service';
import Test from 'models/Test.model';
import useTestCrud, {TTestRunRequest} from './hooks/useTestCrud';

interface IContext {
  onEdit(values: TDraftTest): void;
  onRun(runRequest?: Partial<TTestRunRequest>): void;
  isLoading: boolean;
  isError: boolean;
  test: Test;
  latestTest: Test;
  isLatestVersion: boolean;
  isEditLoading: boolean;
}

export const Context = createContext<IContext>({
  onEdit: noop,
  onRun: noop,
  test: {} as Test,
  latestTest: {} as Test,
  isLoading: false,
  isError: false,
  isLatestVersion: true,
  isEditLoading: false,
});

interface IProps {
  testId: string;
  version?: number;
  children: React.ReactNode;
}

export const useTest = () => useContext(Context);

const TestProvider = ({children, testId, version = 0}: IProps) => {
  const [isVersionModalOpen, setIsVersionModalOpen] = useState(false);
  const [draft, setDraft] = useState<TDraftTest>({});
  const [action, setAction] = useState<'edit' | 'run'>();
  const {runTest, edit, isEditLoading} = useTestCrud();
  const {
    data: test,
    isLoading: isCurrentLoading,
    isError: isCurrentError,
  } = useGetTestVersionByIdQuery({testId, version}, {skip: !version});
  const {data: latestTest, isLoading: isLatestLoading, isError: isLatestError} = useGetTestByIdQuery({testId});

  const isLoading = isLatestLoading || isCurrentLoading;
  const isError = isLatestError || isCurrentError;
  const currentTest = (version ? test : latestTest)!;

  const isLatestVersion = useMemo(
    () => (Boolean(version) && version === latestTest?.version) || currentTest?.version === latestTest?.version,
    [currentTest?.version, latestTest?.version, version]
  );

  const onEdit = useCallback(
    (values: TDraftTest) => {
      if (isLatestVersion) edit(test!, values);
      else {
        setAction('edit');
        setDraft(values);
        setIsVersionModalOpen(true);
      }
    },
    [edit, isLatestVersion, test]
  );

  const onRun = useCallback(
    (request: Partial<TTestRunRequest> = {}) => {
      if (isLatestVersion)
        runTest({
          test: currentTest,
          ...request,
        });
      else {
        setAction('run');
        setIsVersionModalOpen(true);
      }
    },
    [currentTest, isLatestVersion, runTest]
  );

  const onConfirm = useCallback(() => {
    if (action === 'edit') edit(test!, draft);
    else {
      const initialValues = TestService.getInitialValues(test!);
      edit(test!, initialValues);
    }

    setIsVersionModalOpen(false);
  }, [action, draft, edit, test]);

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
          action === 'edit'
            ? 'Editing it will result in a new version that will become the latest.'
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
    <div data-cy="loading-test" />
  );
};

export default TestProvider;
