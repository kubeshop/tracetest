import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {TDraftTransaction, TTransaction} from 'types/Transaction.types';
import {
  useDeleteTransactionByIdMutation,
  useEditTransactionMutation,
  useRunTransactionMutation,
} from 'redux/apis/TraceTest.api';

const useTransactionCrud = () => {
  const navigate = useNavigate();
  const [editTransaction, {isLoading: isTransactionEditLoading}] = useEditTransactionMutation();
  const [runTransactionAction, {isLoading: isLoadingRunTransaction}] = useRunTransactionMutation();
  const [deleteTransactionAction] = useDeleteTransactionByIdMutation();
  const isEditLoading = isTransactionEditLoading || isLoadingRunTransaction;
  const {selectedEnvironment} = useEnvironment();

  const runTransaction = useCallback(
    async (transactionId: string, environmentId = selectedEnvironment?.id) => {
      const run = await runTransactionAction({transactionId, environmentId}).unwrap();

      navigate(`/transaction/${transactionId}/run/${run.id}`);
    },
    [navigate, runTransactionAction, selectedEnvironment?.id]
  );

  const edit = useCallback(
    async (transaction: TTransaction, draft: TDraftTransaction) => {
      const {id: transactionId} = transaction;

      await editTransaction({transactionId, transaction: draft}).unwrap();

      runTransaction(transactionId);
    },
    [editTransaction, runTransaction]
  );

  const deleteTransaction = useCallback(
    (transactionId: string) => {
      deleteTransactionAction({transactionId});

      navigate('/');
    },
    [deleteTransactionAction, navigate]
  );

  return {
    edit,
    runTransaction,
    deleteTransaction,
    isEditLoading,
  };
};

export default useTransactionCrud;
