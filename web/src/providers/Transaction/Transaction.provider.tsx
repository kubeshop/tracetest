import {noop} from 'lodash';
import {createContext, ReactNode, useCallback, useContext, useMemo, useState} from 'react';

import VersionMismatchModal from 'components/VersionMismatchModal';
import Transaction from 'models/Transaction.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useGetTransactionByIdQuery, useGetTransactionVersionByIdQuery} from 'redux/apis/TraceTest.api';
import TransactionService from 'services/Transaction.service';
import {TDraftTransaction} from 'types/Transaction.types';
import useTransactionCrud from './hooks/useTransactionCrud';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';

interface IContext {
  isError: boolean;
  isLoading: boolean;
  isLoadingRun: boolean;
  isEditLoading: boolean;
  onDelete(id: string, name: string): void;
  onEdit(draft: TDraftTransaction): void;
  onRun(runId?: string): void;
  transaction: Transaction;
  latestTransaction: Transaction;
}

export const Context = createContext<IContext>({
  isError: false,
  isLoading: false,
  isLoadingRun: false,
  isEditLoading: false,
  onDelete: noop,
  onRun: noop,
  onEdit: noop,
  transaction: {} as Transaction,
  latestTransaction: {} as Transaction,
});

interface IProps {
  children: ReactNode;
  transactionId: string;
  version?: number;
}

export const useTransaction = () => useContext(Context);

const TransactionProvider = ({children, transactionId, version = 0}: IProps) => {
  const [isVersionModalOpen, setIsVersionModalOpen] = useState(false);
  const [action, setAction] = useState<'edit' | 'run'>();
  const [draft, setDraft] = useState<TDraftTransaction>({});
  const {
    data: latestTransaction,
    isLoading: isLatestLoading,
    isError: isLatestError,
  } = useGetTransactionByIdQuery({transactionId});
  const {deleteTransaction, runTransaction, isEditLoading, edit} = useTransactionCrud();
  const {
    data: transaction,
    isLoading: isCurrentLoading,
    isError: isCurrentError,
  } = useGetTransactionVersionByIdQuery({transactionId, version}, {skip: !version});

  const isLoading = isLatestLoading || isCurrentLoading;
  const isError = isLatestError || isCurrentError;
  const currentTransaction = (version ? transaction : latestTransaction)!;
  const isLatestVersion = useMemo(
    () => Boolean(version) && version === latestTransaction?.version,
    [latestTransaction?.version, version]
  );

  const {onOpen} = useConfirmationModal();
  const {navigate} = useDashboard();

  const onRun = useCallback(
    (runId?: string) => {
      if (isLatestVersion) runTransaction(transaction!, runId);
      else {
        setAction('run');
        setIsVersionModalOpen(true);
      }
    },
    [isLatestVersion, runTransaction, transaction]
  );

  const onDelete = useCallback(
    (id: string, name: string) => {
      function onConfirmation() {
        deleteTransaction(id);
        navigate('/');
      }

      onOpen({
        title: `Are you sure you want to delete “${name}”?`,
        onConfirm: onConfirmation,
      });
    },
    [deleteTransaction, navigate, onOpen]
  );

  const onEdit = useCallback(
    (values: TDraftTransaction) => {
      if (isLatestVersion) edit(transaction!, values);
      else {
        setAction('edit');
        setDraft(values);
        setIsVersionModalOpen(true);
      }
    },
    [edit, isLatestVersion, transaction]
  );

  const onConfirm = useCallback(() => {
    if (action === 'edit') edit(transaction!, draft);
    else {
      const initialValues = TransactionService.getInitialValues(transaction!);
      edit(transaction!, initialValues);
    }

    setIsVersionModalOpen(false);
  }, [action, draft, edit, transaction]);

  const value = useMemo<IContext>(
    () => ({
      isError,
      isLoading,
      isLoadingRun: false,
      onDelete,
      onEdit,
      onRun,
      isEditLoading,
      transaction: currentTransaction!,
      latestTransaction: latestTransaction!,
    }),
    [currentTransaction, isEditLoading, isError, isLoading, latestTransaction, onDelete, onEdit, onRun]
  );

  return currentTransaction && latestTransaction ? (
    <>
      <Context.Provider value={value}>{children}</Context.Provider>
      <VersionMismatchModal
        description={
          action === 'edit'
            ? 'Editing it will result in a new version that will become the latest.'
            : 'Running the transaction will use the latest version of the transaction.'
        }
        currentVersion={currentTransaction.version}
        isOpen={isVersionModalOpen}
        latestVersion={latestTransaction.version}
        okText="Run Transaction"
        onCancel={() => setIsVersionModalOpen(false)}
        onConfirm={onConfirm}
      />
    </>
  ) : (
    <div data-cy="loading-transaction" />
  );
};

export default TransactionProvider;
