import {noop} from 'lodash';
import {createContext, ReactNode, useCallback, useContext, useMemo, useState} from 'react';

import VersionMismatchModal from 'components/VersionMismatchModal';
import TestSuite from 'models/TestSuite.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useGetTestSuiteByIdQuery, useGetTestSuiteVersionByIdQuery} from 'redux/apis/Tracetest';
import TestSuiteService from 'services/TestSuite.service';
import {TDraftTestSuite} from 'types/TestSuite.types';
import useTestSuiteCrud from './hooks/useTestSuiteCrud';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';

interface IContext {
  isError: boolean;
  isLoading: boolean;
  isLoadingRun: boolean;
  isEditLoading: boolean;
  onDelete(id: string, name: string): void;
  onEdit(draft: TDraftTestSuite): void;
  onRun(runId?: string): void;
  testSuite: TestSuite;
  latestTestSuite: TestSuite;
}

export const Context = createContext<IContext>({
  isError: false,
  isLoading: false,
  isLoadingRun: false,
  isEditLoading: false,
  onDelete: noop,
  onRun: noop,
  onEdit: noop,
  testSuite: {} as TestSuite,
  latestTestSuite: {} as TestSuite,
});

interface IProps {
  children: ReactNode;
  testSuiteId: string;
  version?: number;
}

export const useTestSuite = () => useContext(Context);

const TestSuiteProvider = ({children, testSuiteId, version = 0}: IProps) => {
  const [isVersionModalOpen, setIsVersionModalOpen] = useState(false);
  const [action, setAction] = useState<'edit' | 'run'>();
  const [draft, setDraft] = useState<TDraftTestSuite>({});
  const {data: latest, isLoading: isLatestLoading, isError: isLatestError} = useGetTestSuiteByIdQuery({testSuiteId});
  const {deleteTestSuite, runTestSuite, isEditLoading, edit} = useTestSuiteCrud();
  const {
    data: testSuite,
    isLoading: isCurrentLoading,
    isError: isCurrentError,
  } = useGetTestSuiteVersionByIdQuery({testSuiteId, version}, {skip: !version});

  const isLoading = isLatestLoading || isCurrentLoading;
  const isError = isLatestError || isCurrentError;
  const current = (version ? testSuite : latest)!;
  const isLatestVersion = useMemo(() => Boolean(version) && version === latest?.version, [latest?.version, version]);

  const {onOpen} = useConfirmationModal();
  const {navigate} = useDashboard();

  const onRun = useCallback(
    (runId?: string) => {
      if (isLatestVersion) runTestSuite(testSuite!, runId);
      else {
        setAction('run');
        setIsVersionModalOpen(true);
      }
    },
    [isLatestVersion, runTestSuite, testSuite]
  );

  const onDelete = useCallback(
    (id: string, name: string) => {
      function onConfirmation() {
        deleteTestSuite(id);
        navigate('/');
      }

      onOpen({
        title: `Are you sure you want to delete “${name}”?`,
        onConfirm: onConfirmation,
      });
    },
    [deleteTestSuite, navigate, onOpen]
  );

  const onEdit = useCallback(
    (values: TDraftTestSuite) => {
      if (isLatestVersion) edit(testSuite!, values);
      else {
        setAction('edit');
        setDraft(values);
        setIsVersionModalOpen(true);
      }
    },
    [edit, isLatestVersion, testSuite]
  );

  const onConfirm = useCallback(() => {
    if (action === 'edit') edit(testSuite!, draft);
    else {
      const initialValues = TestSuiteService.getInitialValues(testSuite!);
      edit(testSuite!, initialValues);
    }

    setIsVersionModalOpen(false);
  }, [action, draft, edit, testSuite]);

  const value = useMemo<IContext>(
    () => ({
      isError,
      isLoading,
      isLoadingRun: false,
      onDelete,
      onEdit,
      onRun,
      isEditLoading,
      testSuite: current!,
      latestTestSuite: latest!,
    }),
    [isError, isLoading, onDelete, onEdit, onRun, isEditLoading, current, latest]
  );

  return current && latest ? (
    <>
      <Context.Provider value={value}>{children}</Context.Provider>
      <VersionMismatchModal
        description={
          action === 'edit'
            ? 'Editing it will result in a new version that will become the latest.'
            : 'Running the test suite will use the latest version of the test suite.'
        }
        currentVersion={current.version}
        isOpen={isVersionModalOpen}
        latestVersion={latest.version}
        okText="Run Test Suite"
        onCancel={() => setIsVersionModalOpen(false)}
        onConfirm={onConfirm}
      />
    </>
  ) : (
    <div data-cy="loading-testsuite" />
  );
};

export default TestSuiteProvider;
