import {noop} from 'lodash';
import {createContext, ReactNode, useCallback, useContext, useMemo} from 'react';
import {useNavigate} from 'react-router-dom';

import {useGetTransactionByIdQuery} from 'redux/apis/TraceTest.api';
import {TDraftTransaction, TTransaction} from 'types/Transaction.types';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';
import useTransactionCrud from './hooks/useTransactionCrud';

interface IContext {
  isError: boolean;
  isLoading: boolean;
  isLoadingRun: boolean;
  isEditLoading: boolean;
  onDelete(id: string, name: string): void;
  onEdit(draft: TDraftTransaction): void;
  onRun(): void;
  transaction: TTransaction;
}

export const Context = createContext<IContext>({
  isError: false,
  isLoading: false,
  isLoadingRun: false,
  isEditLoading: false,
  onDelete: noop,
  onRun: noop,
  onEdit: noop,
  transaction: {} as TTransaction,
});

interface IProps {
  children: ReactNode;
  transactionId: string;
}

export const useTransaction = () => useContext(Context);

const TransactionProvider = ({children, transactionId}: IProps) => {
  const {data: transaction, isLoading, isError} = useGetTransactionByIdQuery({transactionId});
  const {deleteTransaction, runTransaction, isEditLoading, edit} = useTransactionCrud();
  const {onOpen} = useConfirmationModal();
  const navigate = useNavigate();

  const onRun = useCallback(() => {
    runTransaction(transactionId);
  }, [runTransaction, transactionId]);

  const onDelete = useCallback(
    (id: string, name: string) => {
      function onConfirmation() {
        deleteTransaction(id);
        navigate('/');
      }

      onOpen(`Are you sure you want to delete “${name}”?`, onConfirmation);
    },
    [deleteTransaction, navigate, onOpen]
  );

  const onEdit = useCallback(
    (draft: TDraftTransaction) => {
      edit(transaction!, draft);
    },
    [edit, transaction]
  );

  const value = useMemo<IContext>(
    () => ({
      isError,
      isLoading,
      isLoadingRun: false,
      onDelete,
      onEdit,
      onRun,
      isEditLoading,
      transaction: transaction!,
    }),
    [isEditLoading, isError, isLoading, onDelete, onEdit, onRun, transaction]
  );

  return transaction ? (
    <Context.Provider value={value}>{children}</Context.Provider>
  ) : (
    <div data-cy="loading-transaction" />
  );
};

export default TransactionProvider;
